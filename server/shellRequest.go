package server

import (
	"fmt"
	"github.com/peter-mount/beebrail/server/status"
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
	userData *map[string]interface{} // user data
}

func (r *ShellRequest) Put(k string, v interface{}) {
	(*r.userData)[k] = v
}

func (r *ShellRequest) Get(k string) interface{} {
	return (*r.userData)[k]
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

// convenience method to get the Connection
func (r *ShellRequest) Connection() *status.Connection {
	return (r.Get(KEY_CONNECTION)).(*status.Connection)
}

// convenience method to get the current TableStyle
func (r *ShellRequest) TableStyle() util.TableStyle {
	return (r.Get(KEY_TABLESTYLE)).(util.TableStyle)
}

// convenience method to create a new Table
func (r *ShellRequest) NewTable(title string) *util.Table {
	return &util.Table{
		Title: title,
		Style: r.TableStyle(),
	}
}
