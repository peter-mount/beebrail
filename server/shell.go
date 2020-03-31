package server

import (
	"github.com/peter-mount/beebrail/server/util"
	"github.com/peter-mount/go-telnet/telsh"
	"strings"
)

func (tc *telnetConnection) NewShell(s *Server) *telsh.ShellHandler {
	h := telsh.NewShellHandler()

	// In API mode there is no prompt
	api := tc.config.API
	if api {
		h.Prompt = ""
	} else {
		h.Prompt = tc.config.Shell.Prompt
	}

	h.WelcomeMessage = fixMessage(api, tc.config.Shell.WelcomeMessage)
	h.ExitMessage = fixMessage(api, tc.config.Shell.ExitMessage)

	tc.RegisterShellCommand(h, "crs", s.Crs)
	tc.RegisterShellCommand(h, "depart", s.departures)
	tc.RegisterShellCommand(h, "mode", s.mode)
	tc.RegisterShellCommand(h, "search", s.search)

	return h
}

// fixMessage ensures that a message
func fixMessage(api bool, s string) string {
	// Strip out any CR/LF into unix standard LF
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// If in API mode split it into status messages
	if api {
		var a []string
		util.WriteString(100, s, func(s string) {
			a = append(a, s)
		})
		s = strings.Join(a, "\n") + "\n"
	}

	return strings.ReplaceAll(s, "\n", "\r\n")
}
