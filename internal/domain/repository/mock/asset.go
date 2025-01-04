// Code generated by MockGen. DO NOT EDIT.
// Source: asset.go
//
// Generated by this command:
//
//	mockgen -source=asset.go -destination=mock/asset.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	time "time"

	model "github.com/abyssparanoia/rapid-go/internal/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAsset is a mock of Asset interface.
type MockAsset struct {
	ctrl     *gomock.Controller
	recorder *MockAssetMockRecorder
	isgomock struct{}
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

// GenerateReadPresignedURL mocks base method.
func (m *MockAsset) GenerateReadPresignedURL(ctx context.Context, path string, expires time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateReadPresignedURL", ctx, path, expires)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateReadPresignedURL indicates an expected call of GenerateReadPresignedURL.
func (mr *MockAssetMockRecorder) GenerateReadPresignedURL(ctx, path, expires any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateReadPresignedURL", reflect.TypeOf((*MockAsset)(nil).GenerateReadPresignedURL), ctx, path, expires)
}

// GenerateWritePresignedURL mocks base method.
func (m *MockAsset) GenerateWritePresignedURL(ctx context.Context, contentType model.ContentType, path string, expires time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateWritePresignedURL", ctx, contentType, path, expires)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateWritePresignedURL indicates an expected call of GenerateWritePresignedURL.
func (mr *MockAssetMockRecorder) GenerateWritePresignedURL(ctx, contentType, path, expires any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateWritePresignedURL", reflect.TypeOf((*MockAsset)(nil).GenerateWritePresignedURL), ctx, contentType, path, expires)
}
