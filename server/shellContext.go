package server

import (
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/beebrail/server/util"
)

type ShellContext struct {
	connection     *status.Connection
	responseWriter util.ResponseWriter
	tableStyle     util.TableStyle
	stxEtx         bool
	userData       *map[string]interface{}
	api            bool
}

func (r *ShellContext) Connection() *status.Connection {
	return r.connection
}

func (r *ShellContext) Writer() util.ResponseWriter {
	return r.responseWriter
}

func (r *ShellContext) TableStyle() util.TableStyle {
	return r.tableStyle
}

func (r *ShellContext) SetTableStyle(t util.TableStyle) {
	r.tableStyle = t
}

func (r *ShellContext) IsStxEtx() bool {
	return r.stxEtx
}

func (r *ShellContext) SetStxEtx(f bool) {
	r.stxEtx = f
}

func (r *ShellContext) IsAPI() bool {
	return r.api
}

func (r *ShellContext) Put(k string, v interface{}) {
	(*r.userData)[k] = v
}

func (r *ShellContext) Get(k string) interface{} {
	return (*r.userData)[k]
}
