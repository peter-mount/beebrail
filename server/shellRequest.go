package server

import (
	"fmt"
	"github.com/peter-mount/go-telnet/telsh"
	"io"
	"time"
)

type ShellRequest struct {
	Name   string         // Name of command
	Args   []string       // command arguments
	Stdin  io.ReadCloser  // stdin
	Stdout io.WriteCloser // stdout
	Stderr io.WriteCloser // stderr
}

func (r *ShellRequest) Print(v ...interface{}) *ShellRequest {
	_, _ = fmt.Fprint(r.Stdout, v...)
	return r
}

func (r *ShellRequest) Printf(f string, v ...interface{}) *ShellRequest {
	_, _ = fmt.Fprintf(r.Stdout, f, v...)
	return r
}

func (r *ShellRequest) Println(v ...interface{}) *ShellRequest {
	_, _ = fmt.Fprintln(r.Stdout, v...)
	return r
}

type ShellCommand func(request *ShellRequest) error

func RegisterShellCommand(s *telsh.ShellHandler, name string, c ShellCommand) {
	s.RegisterHandlerFunc(name, func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {
		err := c(&ShellRequest{
			Name:   name,
			Args:   args,
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		})

		// This allows for streams (mainly stdout) to write before the command prompt is written
		// otherwise it appears in the wrong place
		time.Sleep(50 * time.Millisecond)

		return err
	})
}
