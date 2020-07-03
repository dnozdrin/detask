package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		context string
		logPath string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{"test context", args{Test, "stderr"}, Config{Test, "stderr"}},
		{"dev context", args{Dev, "stdout"}, Config{Dev, "stdout"}},
		{"prod context", args{Prod, "file:///dev/null"}, Config{Prod, "file:///dev/null"}},
		{"unknown context", args{mock.Anything, mock.Anything}, Config{Dev, mock.Anything}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewConfig(tt.args.context, tt.args.logPath))
		})
	}
}
