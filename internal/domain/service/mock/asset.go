// Code generated by MockGen. DO NOT EDIT.
// Source: asset.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	model "github.com/abyssparanoia/rapid-go/internal/domain/model"
	service "github.com/abyssparanoia/rapid-go/internal/domain/service"
	gomock "github.com/golang/mock/gomock"
)

// MockAsset is a mock of Asset interface.
type MockAsset struct {
	ctrl     *gomock.Controller
	recorder *MockAssetMockRecorder
}

// MockAssetMockRecorder is the mock recorder for MockAsset.
type MockAssetMockRecorder struct {
	mock *MockAsset
}

// NewMockAsset creates a new mock instance.
func NewMockAsset(ctrl *gomock.Controller) *MockAsset {
	mock := &MockAsset{ctrl: ctrl}
	mock.recorder = &MockAssetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAsset) EXPECT() *MockAssetMockRecorder {
	return m.recorder
}

// BatchSetStaffURLs mocks base method.
func (m *MockAsset) BatchSetStaffURLs(ctx context.Context, staffs model.Staffs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSetStaffURLs", ctx, staffs)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSetStaffURLs indicates an expected call of BatchSetStaffURLs.
func (mr *MockAssetMockRecorder) BatchSetStaffURLs(ctx, staffs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSetStaffURLs", reflect.TypeOf((*MockAsset)(nil).BatchSetStaffURLs), ctx, staffs)
}

// CreatePresignedURL mocks base method.
func (m *MockAsset) CreatePresignedURL(ctx context.Context, assetType model.AssetType, contentType string) (*service.AssetCreatePresignedURLResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePresignedURL", ctx, assetType, contentType)
	ret0, _ := ret[0].(*service.AssetCreatePresignedURLResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePresignedURL indicates an expected call of CreatePresignedURL.
func (mr *MockAssetMockRecorder) CreatePresignedURL(ctx, assetType, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePresignedURL", reflect.TypeOf((*MockAsset)(nil).CreatePresignedURL), ctx, assetType, contentType)
}

// GetWithValidate mocks base method.
func (m *MockAsset) GetWithValidate(ctx context.Context, assetType model.AssetType, assetKey string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithValidate", ctx, assetType, assetKey)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithValidate indicates an expected call of GetWithValidate.
func (mr *MockAssetMockRecorder) GetWithValidate(ctx, assetType, assetKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithValidate", reflect.TypeOf((*MockAsset)(nil).GetWithValidate), ctx, assetType, assetKey)
}