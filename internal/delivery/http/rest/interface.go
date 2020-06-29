package rest

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	"net/http"
	"net/url"
)

type routeAware interface {
	GetURL(name string, params ...string) (*url.URL, error)
	GetIDVar(r *http.Request) (uint, error)
}

type BoardService interface {
	Create(board *m.Board) (*m.Board, error)
	Find() ([]*m.Board, error)
	FindOneById(ID uint) (*m.Board, error)
	Update(board *m.Board) (*m.Board, error)
	Delete(ID uint) error
}

type ColumnService interface {
	Create(board *m.Column) (*m.Column, error)
	Find() ([]*m.Column, error)
	FindOneById(ID uint) (*m.Column, error)
	Update(board *m.Column) (*m.Column, error)
	Delete(ID uint) error
}

type TaskService interface {
	Create(board *m.Task) (*m.Task, error)
	Find(demand services.TaskDemand) ([]*m.Task, error)
	FindOneById(ID uint) (*m.Task, error)
	Update(board *m.Task) (*m.Task, error)
	Delete(ID uint) error
}

type CommentService interface {
	Create(board *m.Comment) (*m.Comment, error)
	Find(demand services.CommentDemand) ([]*m.Comment, error)
	FindOneById(ID uint) (*m.Comment, error)
	Update(board *m.Comment) (*m.Comment, error)
	Delete(ID uint) error
}
