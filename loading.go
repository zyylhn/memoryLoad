package memoryLoad

import (
	"context"
	"github.com/zyylhn/memoryLoad/go-memexec"
)

func Load(ctx context.Context, app *memexec.LoadAppInfo, arg ...string) ([]byte, error) {
	return load(ctx, app, arg...)
}
