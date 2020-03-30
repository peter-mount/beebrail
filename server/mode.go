package server

import (
	"github.com/peter-mount/beebrail/server/util"
	"strings"
)

const (
	KEY_CONNECTION = "connection"
	KEY_TABLESTYLE = "tableStyle"
)

func (s *Server) mode(r *ShellRequest) error {
	for _, m := range r.Args {
		switch strings.ToLower(m) {
		case "plain":
			r.Put(KEY_TABLESTYLE, util.Plain)
		case "sql":
			r.Put(KEY_TABLESTYLE, util.SQL)
		case "bbc":
			r.Put(KEY_TABLESTYLE, util.MODE7)
		default:
			r.Printf("Unsupported mode \"%s\"\r\n", m)
		}
	}
	return nil
}
