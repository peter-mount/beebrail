package server

import (
	"github.com/peter-mount/beebrail/protocol"
	"github.com/peter-mount/golib/kernel"
)

type Server struct {
	protocol *protocol.Protocol
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) Init(k *kernel.Kernel) error {
	service, err := k.AddService(&protocol.Protocol{})
	if err != nil {
		return err
	}
	s.protocol = (service).(*protocol.Protocol)

	return nil
}

func (s *Server) Start() error {

	return nil
}
