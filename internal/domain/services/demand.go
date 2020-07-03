package services

import "github.com/pkg/errors"

// ErrFilterNotAllowed is returned in case of unsupported filter parameters
// are passed
var ErrFilterNotAllowed = errors.Errorf("filter demand is not allowed")

// Demand represents an interface for constraints container
type Demand interface {
	Add(string, uint) error
}

type constraints map[string]uint

var allowedColumnFilter = map[string]struct{}{
	"board": {},
}

// ColumnDemand is a constraints container for tasks
type ColumnDemand constraints

// Add will add allowed filter constraints to the ColumnDemand or will
// return an error if the field / value constraint is not in allowlist
func (cd ColumnDemand) Add(field string, value uint) error {
	if _, ok := allowedColumnFilter[field]; !ok {
		return ErrFilterNotAllowed
	}

	cd[field] = value
	return nil
}

var allowedTaskFilter = map[string]struct{}{
	"board":  {},
	"column": {},
}

// TaskDemand is a constraints container for tasks
type TaskDemand constraints

// Add will add allowed filter constraints to the TaskDemand or will
// return an error if the field / value constraint is not in allowlist
func (td TaskDemand) Add(field string, value uint) error {
	if _, ok := allowedTaskFilter[field]; !ok {
		return ErrFilterNotAllowed
	}

	td[field] = value
	return nil
}

var allowedCommentFilter = map[string]struct{}{
	"task": {},
}

// CommentDemand is a constraints container for comments
type CommentDemand constraints

// Add will add allowed filter constraints to the CommentDemand or will
// return an error if the field / value constraint is not in allowlist
func (cd CommentDemand) Add(field string, value uint) error {
	if _, ok := allowedCommentFilter[field]; !ok {
		return ErrFilterNotAllowed
	}

	cd[field] = value
	return nil
}
