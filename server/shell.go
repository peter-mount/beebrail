package server

import (
	"github.com/peter-mount/go-telnet/telsh"
)

func NewShell(s *Server) *telsh.ShellHandler {
	h := telsh.NewShellHandler()
	h.Prompt = s.config.Shell.Prompt
	h.WelcomeMessage = s.config.Shell.WelcomeMessage
	h.ExitMessage = s.config.Shell.ExitMessage

	RegisterShellCommand(h, "crs", s.Crs)

	return h
}
