package calcium

import (
	"context"
	"testing"

	enginemocks "github.com/projecteru2/core/engine/mocks"
	lockmocks "github.com/projecteru2/core/lock/mocks"
	resourcemocks "github.com/projecteru2/core/resources/mocks"
	storemocks "github.com/projecteru2/core/store/mocks"
	"github.com/projecteru2/core/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPodResource(t *testing.T) {
	c := NewTestCluster()
	ctx := context.Background()
	podname := "testpod"
	nodename := "testnode"
	store := c.store.(*storemocks.Store)
	rmgr := c.rmgr.(*resourcemocks.Manager)
	lock := &lockmocks.DistributedLock{}
	lock.On("Lock", mock.Anything).Return(ctx, nil)
	lock.On("Unlock", mock.Anything).Return(nil)

	// failed by GetNodesByPod
	store.On("GetNodesByPod", mock.Anything, mock.Anything).Return(nil, types.ErrNoETCD).Once()
	ch, err := c.PodResource(ctx, podname)
	assert.Error(t, err)
	store.AssertExpectations(t)
	node := &types.Node{
		NodeMeta: types.NodeMeta{
			Name: nodename,
		},
	}
	store.On("GetNodesByPod", mock.Anything, mock.Anything).Return([]*types.Node{node}, nil)
	store.On("GetNode", mock.Anything, mock.Anything).Return(node, nil)
	store.On("CreateLock", mock.Anything, mock.Anything).Return(lock, nil)

	// failed by ListNodeWorkloads
	store.On("ListNodeWorkloads", mock.Anything, mock.Anything, mock.Anything).Return(nil, types.ErrNoETCD).Once()
	ch, err = c.PodResource(ctx, podname)
	assert.NoError(t, err)
	msg := <-ch
	assert.Equal(t, msg.Name, nodename)
	assert.NotEmpty(t, msg.Diffs)
	store.AssertExpectations(t)
	workloads := []*types.Workload{
		{ResourceArgs: map[string]types.WorkloadResourceArgs{}},
		{ResourceArgs: map[string]types.WorkloadResourceArgs{}},
	}
	store.On("ListNodeWorkloads", mock.Anything, mock.Anything, mock.Anything).Return(workloads, nil)

	// failed by GetNodeResourceInfo
	rmgr.On("GetNodeResourceInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		nil, nil, nil, types.ErrNoETCD).Once()
	ch, err = c.PodResource(ctx, podname)
	msg = <-ch
	assert.NoError(t, err)
	assert.Equal(t, msg.Name, nodename)
	assert.NotEmpty(t, msg.Diffs)
	store.AssertExpectations(t)
	rmgr.On("GetNodeResourceInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		map[string]types.NodeResourceArgs{"test": map[string]interface{}{"abc": 123}},
		map[string]types.NodeResourceArgs{"test": map[string]interface{}{"abc": 123}},
		[]string{},
		nil)

	// success
	ch, err = c.PodResource(ctx, podname)
	msg = <-ch
	assert.NoError(t, err)
	assert.Equal(t, msg.Name, nodename)
	assert.Empty(t, msg.Diffs)
	store.AssertExpectations(t)
}

func TestNodeResource(t *testing.T) {
	c := NewTestCluster()
	ctx := context.Background()
	nodename := "testnode"
	store := c.store.(*storemocks.Store)
	rmgr := c.rmgr.(*resourcemocks.Manager)
	lock := &lockmocks.DistributedLock{}
	store.On("CreateLock", mock.Anything, mock.Anything).Return(lock, nil)
	lock.On("Lock", mock.Anything).Return(ctx, nil)
	lock.On("Unlock", mock.Anything).Return(nil)

	node := &types.Node{
		NodeMeta: types.NodeMeta{
			Name: nodename,
		},
	}
	engine := &enginemocks.API{}
	store.On("GetNode", mock.Anything, mock.Anything).Return(node, nil)
	store.On("CreateLock", mock.Anything, mock.Anything).Return(lock, nil)

	rmgr.On("GetNodeResourceInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		map[string]types.NodeResourceArgs{"test": map[string]interface{}{"abc": 123}},
		map[string]types.NodeResourceArgs{"test": map[string]interface{}{"abc": 123}},
		[]string{},
		nil)

	workloads := []*types.Workload{
		{ResourceArgs: map[string]types.WorkloadResourceArgs{}, Engine: engine},
		{ResourceArgs: map[string]types.WorkloadResourceArgs{}, Engine: engine},
	}
	store.On("ListNodeWorkloads", mock.Anything, mock.Anything, mock.Anything).Return(workloads, nil)
	engine.On("VirtualizationInspect", mock.Anything, mock.Anything).Return(nil, types.ErrNoETCD)

	nr, err := c.NodeResource(ctx, nodename, true)
	assert.NoError(t, err)
	assert.Equal(t, nr.Name, nodename)
	assert.NotEmpty(t, nr.Diffs)
	store.AssertExpectations(t)
}
