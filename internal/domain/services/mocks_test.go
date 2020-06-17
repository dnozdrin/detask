package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/mock"
)

const (
	Anything = "mock.Anything"
)

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

func (mbs *MockedBoardStorage) Update(board *m.Board) (*m.Board, error){
	returnValues := mbs.Called(board)
	return returnValues.Get(0).(*m.Board), returnValues.Error(1)
}

func (mbs *MockedBoardStorage) Delete(ID uint) error {
	returnValues := mbs.Called(ID)
	return returnValues.Error(0)
}

type MockedValidation struct {
	mock.Mock
}

func (mv *MockedValidation) Validate(input interface{}) v.Result {
	returnValues := mv.Called(input)
	return returnValues.Get(0).(v.Result)
}
