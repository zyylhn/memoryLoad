# 内存加载模块

目前只是是linux内存加载，Windows采用创建临时文件的方式实现，后面会支持windows的内存加载功能

使用LoadExecute函数传入程序和参数即可执行并返回结果

```
// LoadExecute 加载给定的二进制程序并返回结果
func LoadExecute(app []byte, args ...string) ([]byte, error) {
	return load("", context.Background(), app, args...)
}

func LoadExecuteWithCtx(ctx context.Context, app []byte, args ...string) ([]byte, error) {
	return load("", ctx, app, args...)
}

//dir用于指定非内存加载的操作系统上临时程序的落盘目录
func Load(dir string, ctx context.Context, app []byte, arg ...string) ([]byte, error) {
	return load(dir, ctx, app, arg...)
}
```

- [x] linux内存加载
- [ ] Windows内存加载（暂时使用落盘执行）
- [ ] 一次性申请过大内存问题