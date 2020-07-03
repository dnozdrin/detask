package services

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCommentDemand_Add(t *testing.T) {
	type args struct {
		field string
		value uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success_task", args{"task", 1}, false},
		{"error", args{mock.Anything, 1}, true},
	}
	demand := make(CommentDemand)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := demand.Add(tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, ErrFilterNotAllowed)
			}
		})
	}
}

func TestTaskDemand_Add(t *testing.T) {
	type args struct {
		field string
		value uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success_board", args{"board", 1}, false},
		{"success_column", args{"column", 1}, false},
		{"error", args{mock.Anything, 1}, true},
	}
	demand := make(TaskDemand)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := demand.Add(tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColumnDemand_Add(t *testing.T) {
	type args struct {
		field string
		value uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success_board", args{"board", 1}, false},
		{"error", args{mock.Anything, 1}, true},
	}
	demand := make(ColumnDemand)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := demand.Add(tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
