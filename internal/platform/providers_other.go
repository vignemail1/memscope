//go:build !windows

// Package platform selects acquisition backends for the current operating system.
package platform

import (
	"context"

	"github.com/vignemail1/memscope/internal/collect"
)

func NewProviders(ctx context.Context) (collect.Providers, error) {
	if err := ctx.Err(); err != nil {
		return collect.Providers{}, err
	}
	return collect.Providers{}, collect.ErrUnsupported
}
