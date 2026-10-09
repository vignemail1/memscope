//go:build windows

// Package platform selects acquisition backends for the current operating system.
package platform

import (
	"context"

	"github.com/vignemail1/memscope/internal/collect"
	windowsbackend "github.com/vignemail1/memscope/internal/platform/windows"
)

func NewProviders(ctx context.Context) (collect.Providers, error) {
	return windowsbackend.NewProviders(ctx)
}
