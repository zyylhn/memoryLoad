package memoryLoad

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/zyylhn/memoryLoad/go-memexec"
	"io"
)

func load(ctx context.Context, app *memexec.LoadAppInfo, args ...string) ([]byte, error) {
	if app.FileName == "" {
		app.AutoDelete = true
	}
	if app.AppBytes == nil && app.AppMaps == nil && app.FileName == "" {
		return nil, fmt.Errorf("when the program content is empty, the file name must be specified")
	}
	exe, err := memexec.New(app)
	if err != nil {
		return nil, err
	}
	defer exe.Close()
	cmd := exe.CommandContext(ctx, args...)
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating StdoutPipe:%v", err))
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating StderrPipe:%v", err))
	}
	defer stderr.Close()

	if err = cmd.Start(); err != nil {
		//fmt.Println("", err)
		return nil, errors.New(fmt.Sprintf("error starting command:%v", err))
	}
	stderrScanner := bufio.NewScanner(stderr)
	var errorMsg string
	go func() {
		for stderrScanner.Scan() {
			errorMsg += stderrScanner.Text() + "\n"
		}
	}()
	var re []byte
	buf := make([]byte, 1024)
	var n int
	if app.MaxResultLen == 0 {
		app.MaxResultLen = 1024 * 1024 * 10
	}
	for {
		n, err = stdout.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		re = append(re, buf[:n]...)
		if len(re) >= app.MaxResultLen {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		//fmt.Println("Error waiting for command:", err)
		return nil, errors.New(fmt.Sprintf("error waiting for command:%v error message'%v'", err, errorMsg))
	}
	if len(errorMsg) > 0 {
		if len(re) == 0 {
			return nil, errors.New("stderr: " + errorMsg)
		} else {
			return re, errors.New("receive stderr:\"" + errorMsg + "\" but there is content in the result")
		}
	}
	return re, nil
}
