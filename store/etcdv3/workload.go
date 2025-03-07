package etcdv3

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/types"
	"github.com/projecteru2/core/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// AddWorkload add a workload
// mainly record its relationship on pod and node
// actually if we already know its node, we will know its pod
// but we still store it
// storage path in etcd is `/workload/:workloadid`
func (m Mercury) AddWorkload(ctx context.Context, workload *types.Workload, processing *types.Processing) error {
	return m.doOpsWorkload(ctx, workload, processing, true)
}

// UpdateWorkload update a workload
func (m *Mercury) UpdateWorkload(ctx context.Context, workload *types.Workload) error {
	return m.doOpsWorkload(ctx, workload, nil, false)
}

// RemoveWorkload remove a workload
// workload id must be in full length
func (m *Mercury) RemoveWorkload(ctx context.Context, workload *types.Workload) error {
	return m.cleanWorkloadData(ctx, workload)
}

// GetWorkload get a workload
// workload if must be in full length, or we can't find it in etcd
// storage path in etcd is `/workload/:workloadid`
func (m *Mercury) GetWorkload(ctx context.Context, ID string) (*types.Workload, error) {
	workloads, err := m.GetWorkloads(ctx, []string{ID})
	if err != nil {
		return nil, err
	}
	return workloads[0], nil
}

// GetWorkloads get many workloads
func (m *Mercury) GetWorkloads(ctx context.Context, ids []string) (workloads []*types.Workload, err error) {
	keys := []string{}
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf(workloadInfoKey, id))
	}

	return m.doGetWorkloads(ctx, keys)
}

// GetWorkloadStatus get workload status
func (m *Mercury) GetWorkloadStatus(ctx context.Context, id string) (*types.StatusMeta, error) {
	workload, err := m.GetWorkload(ctx, id)
	if err != nil {
		return nil, err
	}
	return workload.StatusMeta, nil
}

// SetWorkloadStatus set workload status
func (m *Mercury) SetWorkloadStatus(ctx context.Context, status *types.StatusMeta, ttl int64) error {
	if status.Appname == "" || status.Entrypoint == "" || status.Nodename == "" {
		return types.ErrBadWorkloadStatus
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	statusVal := string(data)
	statusKey := filepath.Join(workloadStatusPrefix, status.Appname, status.Entrypoint, status.Nodename, status.ID)
	workloadKey := fmt.Sprintf(workloadInfoKey, status.ID)
	return m.BindStatus(ctx, workloadKey, statusKey, statusVal, ttl)
}

// ListWorkloads list workloads
func (m *Mercury) ListWorkloads(ctx context.Context, appname, entrypoint, nodename string, limit int64, labels map[string]string) ([]*types.Workload, error) {
	if appname == "" {
		entrypoint = ""
	}
	if entrypoint == "" {
		nodename = ""
	}
	// 这里显式加个 / 来保证 prefix 是唯一的
	key := filepath.Join(workloadDeployPrefix, appname, entrypoint, nodename) + "/"
	resp, err := m.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithLimit(limit))
	if err != nil {
		return nil, err
	}

	workloads := []*types.Workload{}
	for _, ev := range resp.Kvs {
		workload := &types.Workload{}
		if err := json.Unmarshal(ev.Value, workload); err != nil {
			return nil, err
		}
		if utils.LabelsFilter(workload.Labels, labels) {
			workloads = append(workloads, workload)
		}
	}

	return m.bindWorkloadsAdditions(ctx, workloads)
}

// ListNodeWorkloads list workloads belong to one node
func (m *Mercury) ListNodeWorkloads(ctx context.Context, nodename string, labels map[string]string) ([]*types.Workload, error) {
	key := fmt.Sprintf(nodeWorkloadsKey, nodename, "")
	resp, err := m.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	workloads := []*types.Workload{}
	for _, ev := range resp.Kvs {
		workload := &types.Workload{}
		if err := json.Unmarshal(ev.Value, workload); err != nil {
			return nil, err
		}
		if utils.LabelsFilter(workload.Labels, labels) {
			workloads = append(workloads, workload)
		}
	}

	return m.bindWorkloadsAdditions(ctx, workloads)
}

// WorkloadStatusStream watch deployed status
func (m *Mercury) WorkloadStatusStream(ctx context.Context, appname, entrypoint, nodename string, labels map[string]string) chan *types.WorkloadStatus {
	if appname == "" {
		entrypoint = ""
	}
	if entrypoint == "" {
		nodename = ""
	}
	// 显式加个 / 保证 prefix 唯一
	statusKey := filepath.Join(workloadStatusPrefix, appname, entrypoint, nodename) + "/"
	ch := make(chan *types.WorkloadStatus)
	_ = m.pool.Invoke(func() {
		defer func() {
			log.Info("[WorkloadStatusStream] close WorkloadStatus channel")
			close(ch)
		}()

		log.Infof(ctx, "[WorkloadStatusStream] watch on %s", statusKey)
		for resp := range m.Watch(ctx, statusKey, clientv3.WithPrefix()) {
			if resp.Err() != nil {
				if !resp.Canceled {
					log.Errorf(ctx, "[WorkloadStatusStream] watch failed %v", resp.Err())
				}
				return
			}
			for _, ev := range resp.Events {
				_, _, _, ID := parseStatusKey(string(ev.Kv.Key))
				msg := &types.WorkloadStatus{ID: ID, Delete: ev.Type == clientv3.EventTypeDelete}
				workload, err := m.GetWorkload(ctx, ID)
				switch {
				case err != nil:
					msg.Error = err
				case utils.LabelsFilter(workload.Labels, labels):
					log.Debugf(ctx, "[WorkloadStatusStream] workload %s status changed", workload.ID)
					msg.Workload = workload
				default:
					continue
				}
				ch <- msg
			}
		}
	})
	return ch
}

func (m *Mercury) cleanWorkloadData(ctx context.Context, workload *types.Workload) error {
	appname, entrypoint, _, err := utils.ParseWorkloadName(workload.Name)
	if err != nil {
		return err
	}

	keys := []string{
		filepath.Join(workloadStatusPrefix, appname, entrypoint, workload.Nodename, workload.ID), // workload deploy status
		filepath.Join(workloadDeployPrefix, appname, entrypoint, workload.Nodename, workload.ID), // workload deploy status
		fmt.Sprintf(workloadInfoKey, workload.ID),                                                // workload info
		fmt.Sprintf(nodeWorkloadsKey, workload.Nodename, workload.ID),                            // node workloads
	}
	_, err = m.BatchDelete(ctx, keys)
	return err
}

func (m *Mercury) doGetWorkloads(ctx context.Context, keys []string) (workloads []*types.Workload, err error) {
	var kvs []*mvccpb.KeyValue
	if kvs, err = m.GetMulti(ctx, keys); err != nil {
		return
	}

	for _, kv := range kvs {
		workload := &types.Workload{}
		if err = json.Unmarshal(kv.Value, workload); err != nil {
			log.Errorf(ctx, "[doGetWorkloads] failed to unmarshal %v, err: %v", string(kv.Key), err)
			return
		}
		workloads = append(workloads, workload)
	}

	return m.bindWorkloadsAdditions(ctx, workloads)
}

func (m *Mercury) bindWorkloadsAdditions(ctx context.Context, workloads []*types.Workload) ([]*types.Workload, error) {
	nodes := map[string]*types.Node{}
	nodenames := []string{}
	nodenameCache := map[string]struct{}{}
	statusKeys := map[string]string{}
	for _, workload := range workloads {
		appname, entrypoint, _, err := utils.ParseWorkloadName(workload.Name)
		if err != nil {
			return nil, err
		}
		statusKeys[workload.ID] = filepath.Join(workloadStatusPrefix, appname, entrypoint, workload.Nodename, workload.ID)
		if _, ok := nodenameCache[workload.Nodename]; !ok {
			nodenameCache[workload.Nodename] = struct{}{}
			nodenames = append(nodenames, workload.Nodename)
		}
	}
	ns, err := m.GetNodes(ctx, nodenames)
	if err != nil {
		return nil, err
	}
	for _, node := range ns {
		nodes[node.Name] = node
	}

	for index, workload := range workloads {
		if _, ok := nodes[workload.Nodename]; !ok {
			return nil, types.ErrBadMeta
		}
		workloads[index].Engine = nodes[workload.Nodename].Engine
		if _, ok := statusKeys[workload.ID]; !ok {
			continue
		}
		kv, err := m.GetOne(ctx, statusKeys[workload.ID])
		if err != nil {
			continue
		}
		status := &types.StatusMeta{}
		if err := json.Unmarshal(kv.Value, &status); err != nil {
			log.Warnf(ctx, "[bindWorkloadsAdditions] unmarshal %s status data failed %v", workload.ID, err)
			log.Errorf(ctx, "[bindWorkloadsAdditions] status raw: %s", kv.Value)
			continue
		}
		workloads[index].StatusMeta = status
	}
	return workloads, nil
}

func (m *Mercury) doOpsWorkload(ctx context.Context, workload *types.Workload, processing *types.Processing, create bool) error {
	var err error
	appname, entrypoint, _, err := utils.ParseWorkloadName(workload.Name)
	if err != nil {
		return err
	}

	// now everything is ok
	// we use full length id instead
	bytes, err := json.Marshal(workload)
	if err != nil {
		return err
	}
	workloadData := string(bytes)

	data := map[string]string{
		fmt.Sprintf(workloadInfoKey, workload.ID):                                                workloadData,
		fmt.Sprintf(nodeWorkloadsKey, workload.Nodename, workload.ID):                            workloadData,
		filepath.Join(workloadDeployPrefix, appname, entrypoint, workload.Nodename, workload.ID): workloadData,
	}

	var resp *clientv3.TxnResponse
	if create {
		if processing != nil {
			processingKey := m.getProcessingKey(processing)
			err = m.BatchCreateAndDecr(ctx, data, processingKey)
		} else {
			resp, err = m.BatchCreate(ctx, data)
		}
	} else {
		resp, err = m.BatchUpdate(ctx, data)
	}
	if err != nil {
		return err
	}
	if resp != nil && !resp.Succeeded {
		return types.ErrTxnConditionFailed
	}
	return nil
}
