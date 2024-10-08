// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"

	mock "github.com/stretchr/testify/mock"
)

// KubeSeal is an autogenerated mock type for the KubeSeal type
type KubeSeal struct {
	mock.Mock
}

// BuildSecretFile provides a mock function with given fields: secretValues, useRaw
func (_m *KubeSeal) BuildSecretFile(secretValues domain.SecretValues, useRaw bool) {
	_m.Called(secretValues, useRaw)
}

// RawSeal provides a mock function with given fields: secretValues
func (_m *KubeSeal) RawSeal(secretValues domain.SecretValues) {
	_m.Called(secretValues)
}

// Seal provides a mock function with given fields: secret
func (_m *KubeSeal) Seal(secret string) {
	_m.Called(secret)
}

// NewKubeSeal creates a new instance of KubeSeal. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKubeSeal(t interface {
	mock.TestingT
	Cleanup(func())
}) *KubeSeal {
	mock := &KubeSeal{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
