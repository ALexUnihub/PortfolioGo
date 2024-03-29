// Code generated by MockGen. DO NOT EDIT.
// Source: session.go

// Package session is a generated GoMock package.
package session

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionRepo is a mock of SessionRepo interface.
type MockSessionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepoMockRecorder
}

// MockSessionRepoMockRecorder is the mock recorder for MockSessionRepo.
type MockSessionRepoMockRecorder struct {
	mock *MockSessionRepo
}

// NewMockSessionRepo creates a new mock instance.
func NewMockSessionRepo(ctrl *gomock.Controller) *MockSessionRepo {
	mock := &MockSessionRepo{ctrl: ctrl}
	mock.recorder = &MockSessionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepo) EXPECT() *MockSessionRepoMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockSessionRepo) Check(r *http.Request) (*Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", r)
	ret0, _ := ret[0].(*Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockSessionRepoMockRecorder) Check(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockSessionRepo)(nil).Check), r)
}

// Create mocks base method.
func (m *MockSessionRepo) Create(w http.ResponseWriter, userID, userLogin string) (*Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", w, userID, userLogin)
	ret0, _ := ret[0].(*Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionRepoMockRecorder) Create(w, userID, userLogin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionRepo)(nil).Create), w, userID, userLogin)
}
