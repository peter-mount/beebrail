package server

import (
	"github.com/peter-mount/golib/kernel"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	ldbclient "github.com/peter-mount/nre-feeds/ldb/client"
)

type Server struct {
	config    *Config                   // Config file
	refClient refclient.DarwinRefClient // ref api
	ldbClient ldbclient.DarwinLDBClient // ldb api
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) Init(k *kernel.Kernel) error {
	service, err := k.AddService(&Config{})
	if err != nil {
		return err
	}
	s.config = (service).(*Config)

	return nil
}

func (s *Server) PostInit() error {
	s.refClient = refclient.DarwinRefClient{Url: s.config.Services.Reference}
	s.ldbClient = ldbclient.DarwinLDBClient{Url: s.config.Services.LDB}

	return nil
}
