package memoryLoad

import "context"

// LoadExecute 加载给定的二进制程序并返回结果
func LoadExecute(app []byte, args ...string) ([]byte, error) {
	return load("", context.Background(), app, args...)
}

func LoadExecuteWithCtx(ctx context.Context, app []byte, args ...string) ([]byte, error) {
	return load("", ctx, app, args...)
}

func Load(dir string, ctx context.Context, app []byte, arg ...string) ([]byte, error) {
	return load(dir, ctx, app, arg...)
}
