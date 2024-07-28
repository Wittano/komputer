package audio

import (
	"context"
	"io"
)

func Upload(ctx context.Context, r io.ReadCloser) (int, error) {
	defer r.Close()
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
	}

	return 0, nil
}
