package voice

import (
	"context"
	"github.com/wittano/komputer/api"
)

type AudioSearchService interface {
	SearchAudio(ctx context.Context, option AudioSearch, page uint) ([]api.AudioFileInfo, error)
	IsActive() bool
}
