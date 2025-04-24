//go:build !linux
// +build !linux

package memexec

import (
	"os"
	"path/filepath"
	"runtime"
)

//给定文件名和目录并且没有给app的内容就就只返回给定文件的os.File
//给定app内容：
//没给目录名就使用临时目录
//没给文件名就随机文件名
//然后创建文件写入到文件中
//如果给定的文件路径和文件名存在就覆盖掉文件

func open(app *LoadAppInfo) (*os.File, error) {
	if app.Dir == "" {
		app.Dir = os.TempDir()
	}

	info, err := os.Stat(app.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(app.Dir, 0755)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err // 其他错误，例如权限问题
		}
	} else if !info.IsDir() {
		// 路径存在但不是目录
		if err = os.Remove(app.Dir); err != nil {
			return nil, err
		}
		err = os.MkdirAll(app.Dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	if app.AppBytes == nil && app.AppMaps == nil {
		return os.Open(filepath.Join(app.Dir, app.FileName))
	}

	var f *os.File

	if app.FileName == "" {
		app.FileName = "go-memexec-"
		if runtime.GOOS == "windows" {
			app.FileName = "go-memexec-*.exe"
		}
		f, err = os.CreateTemp(app.Dir, app.FileName)
	} else {
		_ = os.Remove(filepath.Join(app.Dir, app.FileName)) // 尝试删掉旧的
		f, err = os.Create(filepath.Join(app.Dir, app.FileName))
	}
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = clean(f, app.AutoDelete)
		}
	}()
	if err = os.Chmod(f.Name(), 0o500); err != nil {
		return nil, err
	}
	if err = app.WriteAppToFile(f); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

func clean(f *os.File, autoDelete bool) error {
	if autoDelete {
		_ = os.Remove(f.Name())
	}
	return f.Close()
}
