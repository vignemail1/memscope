//go:build windows

// Package windows will implement the Windows acquisition backend.
package windows

import (
	"context"

	"github.com/vignemail1/memscope/internal/collect"
)

// NewProviders does not access hardware or install any privileged component.
func NewProviders(ctx context.Context) (collect.Providers, error) {
	if err := ctx.Err(); err != nil {
		return collect.Providers{}, err
	}
	return collect.Providers{}, collect.ErrNotImplemented
}
