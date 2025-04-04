// Code generated by MockGen. DO NOT EDIT.
// Source: tenant.go
//
// Generated by this command:
//
//	mockgen -source=tenant.go -destination=mock/tenant.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/abyssparanoia/rapid-go/internal/domain/model"
	repository "github.com/abyssparanoia/rapid-go/internal/domain/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockTenant is a mock of Tenant interface.
type MockTenant struct {
	ctrl     *gomock.Controller
	recorder *MockTenantMockRecorder
	isgomock struct{}
}

// MockTenantMockRecorder is the mock recorder for MockTenant.
type MockTenantMockRecorder struct {
	mock *MockTenant
}

// NewMockTenant creates a new mock instance.
func NewMockTenant(ctrl *gomock.Controller) *MockTenant {
	mock := &MockTenant{ctrl: ctrl}
	mock.recorder = &MockTenantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTenant) EXPECT() *MockTenantMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockTenant) Count(ctx context.Context, query repository.ListTenantsQuery) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, query)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockTenantMockRecorder) Count(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockTenant)(nil).Count), ctx, query)
}

// Create mocks base method.
func (m *MockTenant) Create(ctx context.Context, tenant *model.Tenant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tenant)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTenantMockRecorder) Create(ctx, tenant any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTenant)(nil).Create), ctx, tenant)
}

// Delete mocks base method.
func (m *MockTenant) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTenantMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTenant)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockTenant) Get(ctx context.Context, query repository.GetTenantQuery) (*model.Tenant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, query)
	ret0, _ := ret[0].(*model.Tenant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTenantMockRecorder) Get(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTenant)(nil).Get), ctx, query)
}

// List mocks base method.
func (m *MockTenant) List(ctx context.Context, query repository.ListTenantsQuery) (model.Tenants, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, query)
	ret0, _ := ret[0].(model.Tenants)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTenantMockRecorder) List(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTenant)(nil).List), ctx, query)
}

// Update mocks base method.
func (m *MockTenant) Update(ctx context.Context, tenant *model.Tenant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tenant)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTenantMockRecorder) Update(ctx, tenant any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTenant)(nil).Update), ctx, tenant)
}
