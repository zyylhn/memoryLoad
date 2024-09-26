//go:build !linux
// +build !linux

package memexec

import (
	"os"
	"runtime"
	"strings"
)

func open(dir string, b []byte) (*os.File, error) {
	pattern := "go-memexec-"
	if runtime.GOOS == "windows" {
		pattern = "go-memexec-*.exe"
	}
R:
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			dir = ""
			goto R
		}
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = clean(f)
		}
	}()
	if err = os.Chmod(f.Name(), 0o500); err != nil {
		return nil, err
	}
	if _, err = f.Write(b); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

func clean(f *os.File) error {
	return os.Remove(f.Name())
}
