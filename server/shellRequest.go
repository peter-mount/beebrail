package server

import (
	"fmt"
	"github.com/peter-mount/beebrail/server/util"
	"io"
)

type ShellRequest struct {
	Name     string                  // Name of command
	Args     []string                // command arguments
	Stdin    io.ReadCloser           // stdin
	Stdout   io.WriteCloser          // stdout
	Stderr   io.WriteCloser          // stderr
	Writer   util.ResponseWriter     // ResponseWriter attached to Stdout
	context  *ShellContext           // Shell context
	userData *map[string]interface{} // user data
}

func (r *ShellRequest) ResultWriter() *util.ResultWriter {
	return util.NewResultWriter(r.Writer)
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

func (r *ShellRequest) Context() *ShellContext {
	return r.context
}

// convenience method to create a new Table
func (r *ShellRequest) NewTable() *util.Table {
	t := &util.Table{Style: r.Context().TableStyle()}
	// Default callbacks
	t.Callback.PageHeader = func(p *util.Pagination, o *util.ResultWriter) error { return nil }
	t.Callback.TableHeader = t.WriteHeader

	return t
}
