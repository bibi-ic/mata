// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bibi-ic/mata/internal/cache (interfaces: MataCache)
//
// Generated by this command:
//
//	mockgen -package mockcache -destination internal/cache/mock/cache.go github.com/bibi-ic/mata/internal/cache MataCache
//
// Package mockcache is a generated GoMock package.
package mockcache

import (
	context "context"
	reflect "reflect"

	models "github.com/bibi-ic/mata/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockMataCache is a mock of MataCache interface.
type MockMataCache struct {
	ctrl     *gomock.Controller
	recorder *MockMataCacheMockRecorder
}

// MockMataCacheMockRecorder is the mock recorder for MockMataCache.
type MockMataCacheMockRecorder struct {
	mock *MockMataCache
}

// NewMockMataCache creates a new mock instance.
func NewMockMataCache(ctrl *gomock.Controller) *MockMataCache {
	mock := &MockMataCache{ctrl: ctrl}
	mock.recorder = &MockMataCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMataCache) EXPECT() *MockMataCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMataCache) Get(arg0 context.Context, arg1 string) (*models.Meta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Meta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMataCacheMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMataCache)(nil).Get), arg0, arg1)
}

// Set mocks base method.
func (m *MockMataCache) Set(arg0 context.Context, arg1 string, arg2 *models.Meta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockMataCacheMockRecorder) Set(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockMataCache)(nil).Set), arg0, arg1, arg2)
}
