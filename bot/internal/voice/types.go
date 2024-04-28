package voice

import (
	"context"
	"github.com/wittano/komputer/api"
	"github.com/wittano/komputer/bot/internal"
)

type AudioSearchService interface {
	AudioFileInfo(ctx context.Context, params SearchParams, page uint) ([]api.AudioFileInfo, error)
	internal.ActiveChecker
}
