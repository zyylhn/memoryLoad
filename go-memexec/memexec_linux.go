//go:build linux
// +build linux

package memexec

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

func open(app *LoadAppInfo) (*os.File, error) {
	if app.AppBytes == nil && app.AppMaps == nil {
		if app.Dir == "" {
			app.Dir = os.TempDir()
		}
		return os.Open(filepath.Join(app.Dir, app.FileName))
	}

	if app.FileName == "" && app.Dir == "" {
		fd, err := unix.MemfdCreate("", unix.MFD_CLOEXEC)
		if err != nil {
			return nil, err
		}
		f := os.NewFile(uintptr(fd), fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd))
		if err = app.WriteAppToFile(f); err != nil {
			_ = f.Close()
			return nil, err
		}
		return f, nil
	} else {
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
		var f *os.File

		if app.FileName == "" {
			app.FileName = "go-memexec-"
			f, err = os.CreateTemp(app.Dir, app.FileName)
		} else {
			_ = os.Remove(filepath.Join(app.Dir, app.FileName)) // 尝试删掉旧的
			f, err = os.Create(filepath.Join(app.Dir, app.FileName))
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
}

func clean(f *os.File, autoDelete bool) error {
	if autoDelete {
		_ = os.Remove(f.Name())
	}
	return f.Close()
}
