package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/mock"
)

const (
	Anything = "mock.Anything"
)

type MockedValidation struct {
	mock.Mock
}

func (mv *MockedValidation) Validate(input interface{}) v.Result {
	returnValues := mv.Called(input)
	return returnValues.Get(0).(v.Result)
}

type MockedBoardStorage struct {
	mock.Mock
}

func (mbs *MockedBoardStorage) Save(board *m.Board) (*m.Board, error) {
	returnValues := mbs.Called(board)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (mbs *MockedBoardStorage) FindById(ID uint) (*m.Board, error) {
	returnValues := mbs.Called(ID)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (mbs *MockedBoardStorage) FindAll() ([]*m.Board, error) {
	returnValues := mbs.Called()
	return returnValues.Get(0).([]*m.Board), returnValues.Error(1)
}

func (mbs *MockedBoardStorage) Update(board *m.Board) (*m.Board, error) {
	returnValues := mbs.Called(board)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (mbs *MockedBoardStorage) Delete(ID uint) error {
	returnValues := mbs.Called(ID)
	return returnValues.Error(0)
}

type MockedColumnStorage struct {
	mock.Mock
}

func (mcs *MockedColumnStorage) Save(board *m.Column) (*m.Column, error) {
	returnValues := mcs.Called(board)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (mcs *MockedColumnStorage) FindById(ID uint) (*m.Column, error) {
	returnValues := mcs.Called(ID)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (mcs *MockedColumnStorage) FindAll() ([]*m.Column, error) {
	returnValues := mcs.Called()
	return returnValues.Get(0).([]*m.Column), returnValues.Error(1)
}

func (mcs *MockedColumnStorage) Update(board *m.Column) (*m.Column, error) {
	returnValues := mcs.Called(board)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (mcs *MockedColumnStorage) Delete(ID uint) error {
	returnValues := mcs.Called(ID)
	return returnValues.Error(0)
}

type MockedTaskStorage struct {
	mock.Mock
}

func (mts *MockedTaskStorage) Save(board *m.Task) (*m.Task, error) {
	returnValues := mts.Called(board)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (mts *MockedTaskStorage) FindById(ID uint) (*m.Task, error) {
	returnValues := mts.Called(ID)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (mts *MockedTaskStorage) FindAll() ([]*m.Task, error) {
	returnValues := mts.Called()
	return returnValues.Get(0).([]*m.Task), returnValues.Error(1)
}

func (mts *MockedTaskStorage) Update(board *m.Task) (*m.Task, error) {
	returnValues := mts.Called(board)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (mts *MockedTaskStorage) Delete(ID uint) error {
	returnValues := mts.Called(ID)
	return returnValues.Error(0)
}

type MockedCommentStorage struct {
	mock.Mock
}

func (mcoms *MockedCommentStorage) Save(board *m.Comment) (*m.Comment, error) {
	returnValues := mcoms.Called(board)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (mcoms *MockedCommentStorage) FindById(ID uint) (*m.Comment, error) {
	returnValues := mcoms.Called(ID)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (mcoms *MockedCommentStorage) FindAll() ([]*m.Comment, error) {
	returnValues := mcoms.Called()
	return returnValues.Get(0).([]*m.Comment), returnValues.Error(1)
}

func (mcoms *MockedCommentStorage) Update(board *m.Comment) (*m.Comment, error) {
	returnValues := mcoms.Called(board)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (mcoms *MockedCommentStorage) Delete(ID uint) error {
	returnValues := mcoms.Called(ID)
	return returnValues.Error(0)
}
