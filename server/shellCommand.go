package server

import (
	"github.com/peter-mount/beebrail/server/util"
	"github.com/peter-mount/go-telnet"
	"github.com/peter-mount/go-telnet/telsh"
	"io"
	"time"
)

type ShellCommand func(request *ShellRequest) error

func (tc *telnetConnection) RegisterShellCommand(s *telsh.ShellHandler, name string, c ShellCommand) {

	produce := func(ctx telnet.Context, name string, args ...string) telsh.Handler {

		handler := func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {
			r := &ShellRequest{
				Name:     name,
				Args:     args,
				Stdin:    stdin,
				Stdout:   stdout,
				Stderr:   stderr,
				userData: ctx.UserData(),
			}

			if tc.config.API {
				r.Writer = util.NewAPIResponseWriter(stdout)
			} else {
				r.Writer = util.NewHumanResponseWriter(stdout)
			}

			err := c(r)

			if err != nil {
				r.Writer.Println(util.ERR_ERROR, err.Error())
			}

			// This allows for streams (mainly stdout) to write before the command prompt is written
			// otherwise it appears in the wrong place
			time.Sleep(50 * time.Millisecond)

			return err
		}

		return PromoteHandlerFunc(ctx, handler, args...)
	}

	producer := telsh.ProducerFunc(produce)

	_ = s.Register(name, producer)
}

type internalPromotedHandlerFunc struct {
	err    error
	fn     telsh.HandlerFunc
	stdin  io.ReadCloser
	stdout io.WriteCloser
	stderr io.WriteCloser

	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser

	args []string

	userData *map[string]interface{}
}

// PromoteHandlerFunc turns a HandlerFunc into a Handler.
func PromoteHandlerFunc(ctx telnet.Context, fn telsh.HandlerFunc, args ...string) telsh.Handler {
	stdin, stdinPipe := io.Pipe()
	stdoutPipe, stdout := io.Pipe()
	stderrPipe, stderr := io.Pipe()

	argsCopy := make([]string, len(args))
	for i, datum := range args {
		argsCopy[i] = datum
	}

	handler := internalPromotedHandlerFunc{
		err: nil,

		fn: fn,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		stdinPipe:  stdinPipe,
		stdoutPipe: stdoutPipe,
		stderrPipe: stderrPipe,

		args:     argsCopy,
		userData: ctx.UserData(),
	}

	return &handler
}

func (handler *internalPromotedHandlerFunc) Run() error {
	if nil != handler.err {
		return handler.err
	}

	handler.err = handler.fn(handler.stdin, handler.stdout, handler.stderr, handler.args...)

	handler.stdin.Close()
	handler.stdout.Close()
	handler.stderr.Close()

	return handler.err
}

func (handler *internalPromotedHandlerFunc) StdinPipe() (io.WriteCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdinPipe, nil
}

func (handler *internalPromotedHandlerFunc) StdoutPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdoutPipe, nil
}

func (handler *internalPromotedHandlerFunc) StderrPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stderrPipe, nil
}
