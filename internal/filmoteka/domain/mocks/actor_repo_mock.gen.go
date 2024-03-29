// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PoorMercymain/filmoteka/internal/filmoteka/domain (interfaces: ActorRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockActorRepository is a mock of ActorRepository interface.
type MockActorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockActorRepositoryMockRecorder
}

// MockActorRepositoryMockRecorder is the mock recorder for MockActorRepository.
type MockActorRepositoryMockRecorder struct {
	mock *MockActorRepository
}

// NewMockActorRepository creates a new mock instance.
func NewMockActorRepository(ctrl *gomock.Controller) *MockActorRepository {
	mock := &MockActorRepository{ctrl: ctrl}
	mock.recorder = &MockActorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActorRepository) EXPECT() *MockActorRepositoryMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockActorRepository) CreateActor(arg0 context.Context, arg1 string, arg2 bool, arg3 time.Time) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockActorRepositoryMockRecorder) CreateActor(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockActorRepository)(nil).CreateActor), arg0, arg1, arg2, arg3)
}

// DeleteActor mocks base method.
func (m *MockActorRepository) DeleteActor(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActor", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActor indicates an expected call of DeleteActor.
func (mr *MockActorRepositoryMockRecorder) DeleteActor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActor", reflect.TypeOf((*MockActorRepository)(nil).DeleteActor), arg0, arg1)
}

// ReadActors mocks base method.
func (m *MockActorRepository) ReadActors(arg0 context.Context, arg1, arg2 int) ([]domain.OutputActor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadActors", arg0, arg1, arg2)
	ret0, _ := ret[0].([]domain.OutputActor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadActors indicates an expected call of ReadActors.
func (mr *MockActorRepositoryMockRecorder) ReadActors(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadActors", reflect.TypeOf((*MockActorRepository)(nil).ReadActors), arg0, arg1, arg2)
}

// UpdateActor mocks base method.
func (m *MockActorRepository) UpdateActor(arg0 context.Context, arg1 int, arg2 string, arg3 *bool, arg4 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockActorRepositoryMockRecorder) UpdateActor(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockActorRepository)(nil).UpdateActor), arg0, arg1, arg2, arg3, arg4)
}
