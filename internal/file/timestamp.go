package file

import (
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"sync"
)

const lockFilePrefix = "/tmp/komputer-"

func CreateLockForService(ctx context.Context, name string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	if IsServiceLocked(name) {
		return
	}

	m := sync.Mutex{}

	m.Lock()
	f, err := os.Create(lockFilePrefix + name)
	if err != nil {
		slog.WarnContext(ctx, "failed lock file for service "+name)
	}
	f.Close()
	m.Unlock()
}

func IsServiceLocked(name string) bool {
	_, err := os.Stat(lockFilePrefix + name)

	return err == nil
}

func RemoveLockForService(ctx context.Context, name string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	if !IsServiceLocked(name) {
		return
	}

	err := os.Remove(lockFilePrefix + name)
	if err != nil {
		slog.ErrorContext(ctx, "failed remove lock file for "+name, err)
	}
}
