// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	enginetypes "github.com/projecteru2/core/engine/types"
	mock "github.com/stretchr/testify/mock"

	resources "github.com/projecteru2/core/resources"

	types "github.com/projecteru2/core/types"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// AddNode provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Manager) AddNode(_a0 context.Context, _a1 string, _a2 types.NodeResourceOpts, _a3 *enginetypes.Info) (map[string]types.NodeResourceArgs, map[string]types.NodeResourceArgs, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, types.NodeResourceOpts, *enginetypes.Info) map[string]types.NodeResourceArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.NodeResourceArgs)
		}
	}

	var r1 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, types.NodeResourceOpts, *enginetypes.Info) map[string]types.NodeResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]types.NodeResourceArgs)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, types.NodeResourceOpts, *enginetypes.Info) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AddPlugins provides a mock function with given fields: _a0
func (_m *Manager) AddPlugins(_a0 ...resources.Plugin) {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Alloc provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Manager) Alloc(_a0 context.Context, _a1 string, _a2 int, _a3 types.WorkloadResourceOpts) ([]types.EngineArgs, []map[string]types.WorkloadResourceArgs, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []types.EngineArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, int, types.WorkloadResourceOpts) []types.EngineArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.EngineArgs)
		}
	}

	var r1 []map[string]types.WorkloadResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, int, types.WorkloadResourceOpts) []map[string]types.WorkloadResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]map[string]types.WorkloadResourceArgs)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int, types.WorkloadResourceOpts) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMetricsDescription provides a mock function with given fields: _a0
func (_m *Manager) GetMetricsDescription(_a0 context.Context) ([]*resources.MetricsDescription, error) {
	ret := _m.Called(_a0)

	var r0 []*resources.MetricsDescription
	if rf, ok := ret.Get(0).(func(context.Context) []*resources.MetricsDescription); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*resources.MetricsDescription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMostIdleNode provides a mock function with given fields: _a0, _a1
func (_m *Manager) GetMostIdleNode(_a0 context.Context, _a1 []string) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, []string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNodeMetrics provides a mock function with given fields: _a0, _a1
func (_m *Manager) GetNodeMetrics(_a0 context.Context, _a1 *types.Node) ([]*resources.Metrics, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*resources.Metrics
	if rf, ok := ret.Get(0).(func(context.Context, *types.Node) []*resources.Metrics); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*resources.Metrics)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.Node) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNodeResourceInfo provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Manager) GetNodeResourceInfo(_a0 context.Context, _a1 string, _a2 []*types.Workload, _a3 bool) (map[string]types.NodeResourceArgs, map[string]types.NodeResourceArgs, []string, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, []*types.Workload, bool) map[string]types.NodeResourceArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.NodeResourceArgs)
		}
	}

	var r1 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, []*types.Workload, bool) map[string]types.NodeResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]types.NodeResourceArgs)
		}
	}

	var r2 []string
	if rf, ok := ret.Get(2).(func(context.Context, string, []*types.Workload, bool) []string); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]string)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, string, []*types.Workload, bool) error); ok {
		r3 = rf(_a0, _a1, _a2, _a3)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// GetNodesDeployCapacity provides a mock function with given fields: _a0, _a1, _a2
func (_m *Manager) GetNodesDeployCapacity(_a0 context.Context, _a1 []string, _a2 types.WorkloadResourceOpts) (map[string]*resources.NodeCapacityInfo, int, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[string]*resources.NodeCapacityInfo
	if rf, ok := ret.Get(0).(func(context.Context, []string, types.WorkloadResourceOpts) map[string]*resources.NodeCapacityInfo); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*resources.NodeCapacityInfo)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, []string, types.WorkloadResourceOpts) int); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, []string, types.WorkloadResourceOpts) error); ok {
		r2 = rf(_a0, _a1, _a2)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPlugins provides a mock function with given fields:
func (_m *Manager) GetPlugins() []resources.Plugin {
	ret := _m.Called()

	var r0 []resources.Plugin
	if rf, ok := ret.Get(0).(func() []resources.Plugin); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resources.Plugin)
		}
	}

	return r0
}

// GetRemapArgs provides a mock function with given fields: _a0, _a1, _a2
func (_m *Manager) GetRemapArgs(_a0 context.Context, _a1 string, _a2 map[string]*types.Workload) (map[string]types.EngineArgs, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[string]types.EngineArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]*types.Workload) map[string]types.EngineArgs); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.EngineArgs)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]*types.Workload) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Realloc provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Manager) Realloc(_a0 context.Context, _a1 string, _a2 map[string]types.WorkloadResourceArgs, _a3 types.WorkloadResourceOpts) (types.EngineArgs, map[string]types.WorkloadResourceArgs, map[string]types.WorkloadResourceArgs, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 types.EngineArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]types.WorkloadResourceArgs, types.WorkloadResourceOpts) types.EngineArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.EngineArgs)
		}
	}

	var r1 map[string]types.WorkloadResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]types.WorkloadResourceArgs, types.WorkloadResourceOpts) map[string]types.WorkloadResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]types.WorkloadResourceArgs)
		}
	}

	var r2 map[string]types.WorkloadResourceArgs
	if rf, ok := ret.Get(2).(func(context.Context, string, map[string]types.WorkloadResourceArgs, types.WorkloadResourceOpts) map[string]types.WorkloadResourceArgs); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(map[string]types.WorkloadResourceArgs)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, string, map[string]types.WorkloadResourceArgs, types.WorkloadResourceOpts) error); ok {
		r3 = rf(_a0, _a1, _a2, _a3)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// RemoveNode provides a mock function with given fields: _a0, _a1
func (_m *Manager) RemoveNode(_a0 context.Context, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollbackAlloc provides a mock function with given fields: _a0, _a1, _a2
func (_m *Manager) RollbackAlloc(_a0 context.Context, _a1 string, _a2 []map[string]types.WorkloadResourceArgs) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []map[string]types.WorkloadResourceArgs) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollbackRealloc provides a mock function with given fields: _a0, _a1, _a2
func (_m *Manager) RollbackRealloc(_a0 context.Context, _a1 string, _a2 map[string]types.WorkloadResourceArgs) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]types.WorkloadResourceArgs) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetNodeResourceCapacity provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *Manager) SetNodeResourceCapacity(_a0 context.Context, _a1 string, _a2 types.NodeResourceOpts, _a3 map[string]types.NodeResourceArgs, _a4 bool, _a5 bool) (map[string]types.NodeResourceArgs, map[string]types.NodeResourceArgs, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, bool, bool) map[string]types.NodeResourceArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.NodeResourceArgs)
		}
	}

	var r1 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, bool, bool) map[string]types.NodeResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]types.NodeResourceArgs)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, bool, bool) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetNodeResourceUsage provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *Manager) SetNodeResourceUsage(_a0 context.Context, _a1 string, _a2 types.NodeResourceOpts, _a3 map[string]types.NodeResourceArgs, _a4 []map[string]types.WorkloadResourceArgs, _a5 bool, _a6 bool) (map[string]types.NodeResourceArgs, map[string]types.NodeResourceArgs, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(0).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, []map[string]types.WorkloadResourceArgs, bool, bool) map[string]types.NodeResourceArgs); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.NodeResourceArgs)
		}
	}

	var r1 map[string]types.NodeResourceArgs
	if rf, ok := ret.Get(1).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, []map[string]types.WorkloadResourceArgs, bool, bool) map[string]types.NodeResourceArgs); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]types.NodeResourceArgs)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, types.NodeResourceOpts, map[string]types.NodeResourceArgs, []map[string]types.WorkloadResourceArgs, bool, bool) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewManager creates a new instance of Manager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewManager(t mockConstructorTestingTNewManager) *Manager {
	mock := &Manager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
