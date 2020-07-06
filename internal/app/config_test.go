// +build unit

package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		context        string
		logPath        string
		allowedOrigins string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"test_context",
			args{Test, "stderr", ""},
			Config{Test, "stderr", []string{""}},
		},
		{
			"dev_ontext",
			args{Dev, "stdout", "http://localhost:8080"},
			Config{Dev, "stdout", []string{"http://localhost:8080"}}},
		{
			"prod_context",
			args{Prod, "file:///dev/null", "http://localhost:8080,http://localhost:80"},
			Config{Prod, "file:///dev/null", []string{"http://localhost:8080", "http://localhost:80"}},
		},		{
			"whitespaces_origings",
			args{Dev, "stderr", "http://localhost:8080, http://localhost:80 "},
			Config{Dev, "stderr", []string{"http://localhost:8080", "http://localhost:80"}},
		},
		{
			"unknown_context",
			args{mock.Anything, mock.Anything, ""},
			Config{Dev, mock.Anything, []string{""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewConfig(tt.args.context, tt.args.logPath, tt.args.allowedOrigins))
		})
	}
}
