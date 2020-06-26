package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/mock"
)

type MockedValidation struct {
	mock.Mock
}

func (val *MockedValidation) Validate(input interface{}) *v.Errors {
	returnValues := val.Called(input)
	return returnValues.Get(0).(*v.Errors)
}

var _ BoardStorage = new(MockedBoardStorage)

type MockedBoardStorage struct {
	mock.Mock
}

func (bs *MockedBoardStorage) SaveWithDefaultColumn(board *m.Board) (*m.Board, error) {
	returnValues := bs.Called(board)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (bs *MockedBoardStorage) FindOneById(ID uint) (*m.Board, error) {
	returnValues := bs.Called(ID)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (bs *MockedBoardStorage) Find() ([]*m.Board, error) {
	returnValues := bs.Called()
	return returnValues.Get(0).([]*m.Board), returnValues.Error(1)
}

func (bs *MockedBoardStorage) Update(board *m.Board) (*m.Board, error) {
	returnValues := bs.Called(board)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (bs *MockedBoardStorage) Delete(ID uint) error {
	returnValues := bs.Called(ID)
	return returnValues.Error(0)
}

var _ ColumnStorage = new(MockedColumnStorage)

type MockedColumnStorage struct {
	mock.Mock
}

func (cs *MockedColumnStorage) Save(column *m.Column) (*m.Column, error) {
	returnValues := cs.Called(column)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (cs *MockedColumnStorage) FindOneById(ID uint) (*m.Column, error) {
	returnValues := cs.Called(ID)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (cs *MockedColumnStorage) Find() ([]*m.Column, error) {
	returnValues := cs.Called()
	return returnValues.Get(0).([]*m.Column), returnValues.Error(1)
}

func (cs *MockedColumnStorage) Update(column *m.Column) (*m.Column, error) {
	returnValues := cs.Called(column)
	return returnValues.Get(0).(*m.Column), returnValues.Error(1)
}

func (cs *MockedColumnStorage) Delete(ID uint) error {
	returnValues := cs.Called(ID)
	return returnValues.Error(0)
}

var _ TaskStorage = new(MockedTaskStorage)

type MockedTaskStorage struct {
	mock.Mock
}

func (ts *MockedTaskStorage) Save(task *m.Task) (*m.Task, error) {
	returnValues := ts.Called(task)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (ts *MockedTaskStorage) FindOneById(ID uint) (*m.Task, error) {
	returnValues := ts.Called(ID)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (ts *MockedTaskStorage) Find(demand TaskDemand) ([]*m.Task, error) {
	returnValues := ts.Called(demand)
	return returnValues.Get(0).([]*m.Task), returnValues.Error(1)
}

func (ts *MockedTaskStorage) Update(task *m.Task) (*m.Task, error) {
	returnValues := ts.Called(task)
	return returnValues.Get(0).(*m.Task), returnValues.Error(1)
}

func (ts *MockedTaskStorage) Delete(ID uint) error {
	returnValues := ts.Called(ID)
	return returnValues.Error(0)
}

var _ CommentStorage = new(MockedCommentStorage)

type MockedCommentStorage struct {
	mock.Mock
}

func (coms *MockedCommentStorage) Save(comment *m.Comment) (*m.Comment, error) {
	returnValues := coms.Called(comment)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (coms *MockedCommentStorage) FindOneById(ID uint) (*m.Comment, error) {
	returnValues := coms.Called(ID)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (coms *MockedCommentStorage) Find(demand CommentDemand) ([]*m.Comment, error) {
	returnValues := coms.Called(demand)
	return returnValues.Get(0).([]*m.Comment), returnValues.Error(1)
}

func (coms *MockedCommentStorage) Update(comment *m.Comment) (*m.Comment, error) {
	returnValues := coms.Called(comment)
	return returnValues.Get(0).(*m.Comment), returnValues.Error(1)
}

func (coms *MockedCommentStorage) Delete(ID uint) error {
	returnValues := coms.Called(ID)
	return returnValues.Error(0)
}
