// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kripsy/shortener/internal/app/handlers (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/kripsy/shortener/internal/app/models"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRepository) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// CreateOrGetBatchFromStorage mocks base method.
func (m *MockRepository) CreateOrGetBatchFromStorage(arg0 context.Context, arg1 *models.BatchURL) (*models.BatchURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrGetBatchFromStorage", arg0, arg1)
	ret0, _ := ret[0].(*models.BatchURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrGetBatchFromStorage indicates an expected call of CreateOrGetBatchFromStorage.
func (mr *MockRepositoryMockRecorder) CreateOrGetBatchFromStorage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrGetBatchFromStorage", reflect.TypeOf((*MockRepository)(nil).CreateOrGetBatchFromStorage), arg0, arg1)
}

// CreateOrGetFromStorage mocks base method.
func (m *MockRepository) CreateOrGetFromStorage(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrGetFromStorage", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrGetFromStorage indicates an expected call of CreateOrGetFromStorage.
func (mr *MockRepositoryMockRecorder) CreateOrGetFromStorage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrGetFromStorage", reflect.TypeOf((*MockRepository)(nil).CreateOrGetFromStorage), arg0, arg1)
}

// GetOriginalURLFromStorage mocks base method.
func (m *MockRepository) GetOriginalURLFromStorage(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURLFromStorage", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURLFromStorage indicates an expected call of GetOriginalURLFromStorage.
func (mr *MockRepositoryMockRecorder) GetOriginalURLFromStorage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURLFromStorage", reflect.TypeOf((*MockRepository)(nil).GetOriginalURLFromStorage), arg0, arg1)
}

// Ping mocks base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}
