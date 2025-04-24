package memoryLoad

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/zyylhn/memoryLoad/go-memexec"
	"log"
	"testing"
	"time"
)

//go:embed cmd/app/zscan
var expapp []byte

// 测试执行后自动删除
func TestLoadExecute(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*50)
	re, err := Load(ctx, &memexec.LoadAppInfo{
		FileName:   "",
		Dir:        "",
		AppBytes:   expapp,
		AutoDelete: false,
	}, "ps", "-H", "10.224.1.10")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(re))
	c()
}

// 测试指定目录
func TestLoadExecute2(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*50)
	re, err := Load(ctx, &memexec.LoadAppInfo{
		FileName:   "",
		Dir:        "/tmp/",
		AppBytes:   expapp,
		AutoDelete: false,
	}, "ps", "-H", "10.224.1.10", "-T", "1", "-p", "1-5000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(re))
	c()
}

// 测试指定文件指定目录自动删除是否生效
func TestLoadExecute3(t *testing.T) {
	ctx := context.Background()
	re, err := Load(ctx, &memexec.LoadAppInfo{
		FileName:   "zscan",
		Dir:        "/tmp",
		AppBytes:   expapp,
		AutoDelete: false,
	}, "ps", "-H", "10.224.1.10")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(re))
}

// 测试直接给定文件路径，不给文件内容
func TestLoadExecute4(t *testing.T) {
	ctx := context.Background()
	re, err := Load(ctx, &memexec.LoadAppInfo{
		FileName:   "zscan",
		Dir:        "/tmp",
		AppBytes:   nil,
		AutoDelete: false,
	}, "ps", "-H", "10.224.1.10")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(re))
}
