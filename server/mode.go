package server

import (
	"github.com/peter-mount/beebrail/server/util"
	"strings"
)

func (s *Server) mode(r *ShellRequest) error {
	var a []string
	ctx := r.Context()
	for _, m := range r.Args {
		switch strings.ToLower(m) {
		case "plain":
			ctx.SetTableStyle(util.Plain)
		case "sql":
			ctx.SetTableStyle(util.SQL)
		case "bbc":
			ctx.SetTableStyle(util.MODE7)
		case "api":
			ctx.SetStxEtx(!ctx.IsStxEtx())
			// Prefix ! to indicate it's off
			if !ctx.IsStxEtx() {
				m = "!" + m
			}
		default:
			r.Printf("Unsupported mode \"%s\"\r\n", m)
			continue
		}

		a = append(a, m)
	}

	if ctx.IsAPI() {
		r.Writer.Println(util.ERR_OK, strings.Join(a, "\n"))
	}

	return nil
}
