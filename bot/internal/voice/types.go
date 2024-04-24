package voice

import (
	"context"
	"github.com/wittano/komputer/api"
	"github.com/wittano/komputer/db"
)

type AudioSearchService interface {
	SearchAudio(ctx context.Context, option db.AudioSearch, page uint) ([]api.AudioFileInfo, error)
}
