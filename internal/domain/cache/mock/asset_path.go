// Code generated by MockGen. DO NOT EDIT.
// Source: asset_path.go
//
// Generated by this command:
//
//	mockgen -source=asset_path.go -destination=mock/asset_path.go -package=mock_cache
//

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	reflect "reflect"

	model "github.com/abyssparanoia/rapid-go/internal/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAssetPath is a mock of AssetPath interface.
type MockAssetPath struct {
	ctrl     *gomock.Controller
	recorder *MockAssetPathMockRecorder
	isgomock struct{}
}

// MockAssetPathMockRecorder is the mock recorder for MockAssetPath.
type MockAssetPathMockRecorder struct {
	mock *MockAssetPath
}

// NewMockAssetPath creates a new mock instance.
func NewMockAssetPath(ctrl *gomock.Controller) *MockAssetPath {
	mock := &MockAssetPath{ctrl: ctrl}
	mock.recorder = &MockAssetPathMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssetPath) EXPECT() *MockAssetPathMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockAssetPath) Clear(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockAssetPathMockRecorder) Clear(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockAssetPath)(nil).Clear), ctx, id)
}

// Get mocks base method.
func (m *MockAssetPath) Get(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAssetPathMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAssetPath)(nil).Get), ctx, id)
}

// Set mocks base method.
func (m *MockAssetPath) Set(ctx context.Context, asset *model.Asset) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, asset)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockAssetPathMockRecorder) Set(ctx, asset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockAssetPath)(nil).Set), ctx, asset)
}
