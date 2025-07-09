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
			if len(errorMsg) > app.MaxResultLen {
				errorMsg = errorMsg[:app.MaxResultLen]
				break
			}
		}
	}()
	var re []byte
	buf := make([]byte, 1024)
	rb := NewRingBuffer(app.MaxResultLen)
	var n int
	for {
		n, err = stdout.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rb.Write(buf[:n])
	}
	re = rb.Bytes()
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

type RingBuffer struct {
	data  []byte
	size  int
	pos   int
	full  bool
	unLim bool // 是否不限制长度
}

func NewRingBuffer(size int) *RingBuffer {
	if size <= 0 {
		return &RingBuffer{
			data:  make([]byte, 0, 4096),
			unLim: true,
		}
	}
	return &RingBuffer{
		data: make([]byte, size),
		size: size,
	}
}

func (r *RingBuffer) Write(p []byte) {
	if r.unLim {
		r.data = append(r.data, p...)
		return
	}
	for _, b := range p {
		r.data[r.pos] = b
		r.pos = (r.pos + 1) % r.size
		if r.pos == 0 {
			r.full = true
		}
	}
}

func (r *RingBuffer) Bytes() []byte {
	if r.unLim {
		return r.data
	}
	if !r.full {
		return r.data[:r.pos]
	}
	return append(r.data[r.pos:], r.data[:r.pos]...)
}
