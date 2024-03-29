// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PoorMercymain/filmoteka/internal/filmoteka/domain (interfaces: FilmRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockFilmRepository is a mock of FilmRepository interface.
type MockFilmRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFilmRepositoryMockRecorder
}

// MockFilmRepositoryMockRecorder is the mock recorder for MockFilmRepository.
type MockFilmRepositoryMockRecorder struct {
	mock *MockFilmRepository
}

// NewMockFilmRepository creates a new mock instance.
func NewMockFilmRepository(ctrl *gomock.Controller) *MockFilmRepository {
	mock := &MockFilmRepository{ctrl: ctrl}
	mock.recorder = &MockFilmRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmRepository) EXPECT() *MockFilmRepositoryMockRecorder {
	return m.recorder
}

// CreateFilm mocks base method.
func (m *MockFilmRepository) CreateFilm(arg0 context.Context, arg1, arg2 string, arg3 time.Time, arg4 float32, arg5 []int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFilm", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFilm indicates an expected call of CreateFilm.
func (mr *MockFilmRepositoryMockRecorder) CreateFilm(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFilm", reflect.TypeOf((*MockFilmRepository)(nil).CreateFilm), arg0, arg1, arg2, arg3, arg4, arg5)
}

// DeleteFilm mocks base method.
func (m *MockFilmRepository) DeleteFilm(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFilm", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFilm indicates an expected call of DeleteFilm.
func (mr *MockFilmRepositoryMockRecorder) DeleteFilm(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFilm", reflect.TypeOf((*MockFilmRepository)(nil).DeleteFilm), arg0, arg1)
}

// FindFilms mocks base method.
func (m *MockFilmRepository) FindFilms(arg0 context.Context, arg1, arg2 string, arg3, arg4 int) ([]domain.OutputFilm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindFilms", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]domain.OutputFilm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFilms indicates an expected call of FindFilms.
func (mr *MockFilmRepositoryMockRecorder) FindFilms(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFilms", reflect.TypeOf((*MockFilmRepository)(nil).FindFilms), arg0, arg1, arg2, arg3, arg4)
}

// ReadFilms mocks base method.
func (m *MockFilmRepository) ReadFilms(arg0 context.Context, arg1, arg2 string, arg3, arg4 int) ([]domain.OutputFilm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFilms", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]domain.OutputFilm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFilms indicates an expected call of ReadFilms.
func (mr *MockFilmRepositoryMockRecorder) ReadFilms(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFilms", reflect.TypeOf((*MockFilmRepository)(nil).ReadFilms), arg0, arg1, arg2, arg3, arg4)
}

// UpdateFilm mocks base method.
func (m *MockFilmRepository) UpdateFilm(arg0 context.Context, arg1 int, arg2, arg3 string, arg4 time.Time, arg5 *float32, arg6 []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFilm", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFilm indicates an expected call of UpdateFilm.
func (mr *MockFilmRepositoryMockRecorder) UpdateFilm(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFilm", reflect.TypeOf((*MockFilmRepository)(nil).UpdateFilm), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}
