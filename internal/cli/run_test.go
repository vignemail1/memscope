package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		code int
		out  string
		err  string
	}{
		{"default help", nil, 0, "Usage:", ""},
		{"help", []string{"--help"}, 0, "Usage:", ""},
		{"short help", []string{"-h"}, 0, "Usage:", ""},
		{"version", []string{"--version"}, 0, "memscope test\n", ""},
		{"unimplemented", []string{"inspect"}, 3, "", "not implemented"},
		{"unknown", []string{"invalid"}, 2, "", "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := Run(tt.args, &stdout, &stderr, "test"); got != tt.code {
				t.Fatalf("code = %d, want %d", got, tt.code)
			}
			if !strings.Contains(stdout.String(), tt.out) || (tt.out == "" && stdout.Len() != 0) {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.err) || (tt.err == "" && stderr.Len() != 0) {
				t.Fatalf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestRunWriteFailure(t *testing.T) {
	if got := Run(nil, failingWriter{}, &bytes.Buffer{}, "test"); got != 1 {
		t.Fatalf("code = %d, want 1", got)
	}
}
