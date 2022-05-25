// Code generated by MockGen. DO NOT EDIT.
// Source: comment.go

// Package comments is a generated GoMock package.
package comments

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mgo_v2 "gopkg.in/mgo.v2"
)

// MockCommentsRepo is a mock of CommentsRepo interface.
type MockCommentsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCommentsRepoMockRecorder
}

// MockCommentsRepoMockRecorder is the mock recorder for MockCommentsRepo.
type MockCommentsRepoMockRecorder struct {
	mock *MockCommentsRepo
}

// NewMockCommentsRepo creates a new mock instance.
func NewMockCommentsRepo(ctrl *gomock.Controller) *MockCommentsRepo {
	mock := &MockCommentsRepo{ctrl: ctrl}
	mock.recorder = &MockCommentsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentsRepo) EXPECT() *MockCommentsRepoMockRecorder {
	return m.recorder
}

// AddStaticComment mocks base method.
func (m *MockCommentsRepo) AddStaticComment(session *mgo_v2.Session, postID string) (*Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddStaticComment", session, postID)
	ret0, _ := ret[0].(*Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddStaticComment indicates an expected call of AddStaticComment.
func (mr *MockCommentsRepoMockRecorder) AddStaticComment(session, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStaticComment", reflect.TypeOf((*MockCommentsRepo)(nil).AddStaticComment), session, postID)
}

// CreateComment mocks base method.
func (m *MockCommentsRepo) CreateComment(session *mgo_v2.Session, postID, data, userID, userLogin string) (*Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", session, postID, data, userID, userLogin)
	ret0, _ := ret[0].(*Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockCommentsRepoMockRecorder) CreateComment(session, postID, data, userID, userLogin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockCommentsRepo)(nil).CreateComment), session, postID, data, userID, userLogin)
}

// DeleteCommentFromRepo mocks base method.
func (m *MockCommentsRepo) DeleteCommentFromRepo(session *mgo_v2.Session, postID, commID, userID, userLogin string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommentFromRepo", session, postID, commID, userID, userLogin)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommentFromRepo indicates an expected call of DeleteCommentFromRepo.
func (mr *MockCommentsRepoMockRecorder) DeleteCommentFromRepo(session, postID, commID, userID, userLogin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommentFromRepo", reflect.TypeOf((*MockCommentsRepo)(nil).DeleteCommentFromRepo), session, postID, commID, userID, userLogin)
}

// GetAllComments mocks base method.
func (m *MockCommentsRepo) GetAllComments(session *mgo_v2.Session, postID string) ([]*Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllComments", session, postID)
	ret0, _ := ret[0].([]*Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllComments indicates an expected call of GetAllComments.
func (mr *MockCommentsRepoMockRecorder) GetAllComments(session, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllComments", reflect.TypeOf((*MockCommentsRepo)(nil).GetAllComments), session, postID)
}