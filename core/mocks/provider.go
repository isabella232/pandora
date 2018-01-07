// Code generated by mockery v1.0.0
package coremock

import context "context"
import core "github.com/yandex/pandora/core"
import mock "github.com/stretchr/testify/mock"

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// Acquire provides a mock function with given fields:
func (_m *Provider) Acquire() (core.Ammo, bool) {
	ret := _m.Called()

	if len(ret) == 1 {
		if rf, ok := ret.Get(0).(func() (core.Ammo, bool)); ok {
			return rf()
		}
	}

	var r0 core.Ammo
	if rf, ok := ret.Get(0).(func() core.Ammo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Ammo)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Release provides a mock function with given fields: ammo
func (_m *Provider) Release(ammo core.Ammo) {
	_m.Called(ammo)
}

// Run provides a mock function with given fields: ctx, deps
func (_m *Provider) Run(ctx context.Context, deps core.ProviderDeps) error {
	ret := _m.Called(ctx, deps)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.ProviderDeps) error); ok {
		r0 = rf(ctx, deps)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
