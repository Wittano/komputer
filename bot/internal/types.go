package internal

import "context"

type ActiveChecker interface {
	Active(ctx context.Context) bool
}
