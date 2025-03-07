// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	clientv3 "go.etcd.io/etcd/client/v3"

	mock "github.com/stretchr/testify/mock"
)

// ETCDClientV3 is an autogenerated mock type for the ETCDClientV3 type
type ETCDClientV3 struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *ETCDClientV3) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Compact provides a mock function with given fields: ctx, rev, opts
func (_m *ETCDClientV3) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, rev)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clientv3.CompactResponse
	if rf, ok := ret.Get(0).(func(context.Context, int64, ...clientv3.CompactOption) *clientv3.CompactResponse); ok {
		r0 = rf(ctx, rev, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.CompactResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, ...clientv3.CompactOption) error); ok {
		r1 = rf(ctx, rev, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, key, opts
func (_m *ETCDClientV3) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clientv3.DeleteResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, ...clientv3.OpOption) *clientv3.DeleteResponse); ok {
		r0 = rf(ctx, key, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.DeleteResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...clientv3.OpOption) error); ok {
		r1 = rf(ctx, key, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Do provides a mock function with given fields: ctx, op
func (_m *ETCDClientV3) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	ret := _m.Called(ctx, op)

	var r0 clientv3.OpResponse
	if rf, ok := ret.Get(0).(func(context.Context, clientv3.Op) clientv3.OpResponse); ok {
		r0 = rf(ctx, op)
	} else {
		r0 = ret.Get(0).(clientv3.OpResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, clientv3.Op) error); ok {
		r1 = rf(ctx, op)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, key, opts
func (_m *ETCDClientV3) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clientv3.GetResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, ...clientv3.OpOption) *clientv3.GetResponse); ok {
		r0 = rf(ctx, key, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.GetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...clientv3.OpOption) error); ok {
		r1 = rf(ctx, key, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Grant provides a mock function with given fields: ctx, ttl
func (_m *ETCDClientV3) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	ret := _m.Called(ctx, ttl)

	var r0 *clientv3.LeaseGrantResponse
	if rf, ok := ret.Get(0).(func(context.Context, int64) *clientv3.LeaseGrantResponse); ok {
		r0 = rf(ctx, ttl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.LeaseGrantResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, ttl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeepAlive provides a mock function with given fields: ctx, id
func (_m *ETCDClientV3) KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 <-chan *clientv3.LeaseKeepAliveResponse
	if rf, ok := ret.Get(0).(func(context.Context, clientv3.LeaseID) <-chan *clientv3.LeaseKeepAliveResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *clientv3.LeaseKeepAliveResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, clientv3.LeaseID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeepAliveOnce provides a mock function with given fields: ctx, id
func (_m *ETCDClientV3) KeepAliveOnce(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *clientv3.LeaseKeepAliveResponse
	if rf, ok := ret.Get(0).(func(context.Context, clientv3.LeaseID) *clientv3.LeaseKeepAliveResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.LeaseKeepAliveResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, clientv3.LeaseID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Leases provides a mock function with given fields: ctx
func (_m *ETCDClientV3) Leases(ctx context.Context) (*clientv3.LeaseLeasesResponse, error) {
	ret := _m.Called(ctx)

	var r0 *clientv3.LeaseLeasesResponse
	if rf, ok := ret.Get(0).(func(context.Context) *clientv3.LeaseLeasesResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.LeaseLeasesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put provides a mock function with given fields: ctx, key, val, opts
func (_m *ETCDClientV3) Put(ctx context.Context, key string, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, key, val)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clientv3.PutResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...clientv3.OpOption) *clientv3.PutResponse); ok {
		r0 = rf(ctx, key, val, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.PutResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...clientv3.OpOption) error); ok {
		r1 = rf(ctx, key, val, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RequestProgress provides a mock function with given fields: ctx
func (_m *ETCDClientV3) RequestProgress(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Revoke provides a mock function with given fields: ctx, id
func (_m *ETCDClientV3) Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *clientv3.LeaseRevokeResponse
	if rf, ok := ret.Get(0).(func(context.Context, clientv3.LeaseID) *clientv3.LeaseRevokeResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.LeaseRevokeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, clientv3.LeaseID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TimeToLive provides a mock function with given fields: ctx, id, opts
func (_m *ETCDClientV3) TimeToLive(ctx context.Context, id clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clientv3.LeaseTimeToLiveResponse
	if rf, ok := ret.Get(0).(func(context.Context, clientv3.LeaseID, ...clientv3.LeaseOption) *clientv3.LeaseTimeToLiveResponse); ok {
		r0 = rf(ctx, id, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clientv3.LeaseTimeToLiveResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, clientv3.LeaseID, ...clientv3.LeaseOption) error); ok {
		r1 = rf(ctx, id, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Txn provides a mock function with given fields: ctx
func (_m *ETCDClientV3) Txn(ctx context.Context) clientv3.Txn {
	ret := _m.Called(ctx)

	var r0 clientv3.Txn
	if rf, ok := ret.Get(0).(func(context.Context) clientv3.Txn); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(clientv3.Txn)
		}
	}

	return r0
}

// Watch provides a mock function with given fields: ctx, key, opts
func (_m *ETCDClientV3) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 clientv3.WatchChan
	if rf, ok := ret.Get(0).(func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan); ok {
		r0 = rf(ctx, key, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(clientv3.WatchChan)
		}
	}

	return r0
}

type mockConstructorTestingTNewETCDClientV3 interface {
	mock.TestingT
	Cleanup(func())
}

// NewETCDClientV3 creates a new instance of ETCDClientV3. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewETCDClientV3(t mockConstructorTestingTNewETCDClientV3) *ETCDClientV3 {
	mock := &ETCDClientV3{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
