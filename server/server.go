package server

import (
	"bufio"
	"errors"
	"github.com/peter-mount/golib/kernel"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	ldbclient "github.com/peter-mount/nre-feeds/ldb/client"
	"io"
	"log"
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

// ReadPacket reads a packet from a Reader
func (s *Server) ReadPacket(in *bufio.Reader, out io.Writer) (Packet, error) {
	cmd := Packet{out: out}
	err := cmd.Read(in)
	return cmd, err
}

// ProcessPacket processes a Packet
func (s *Server) ProcessPacket(cmd Packet) error {
	if cmd.out == nil {
		return errors.New("no response stream")
	}

	var resp *Packet
	switch cmd.Command {
	case 'C':
		resp = s.crs(cmd)
	case 'D':
		resp = s.departures(cmd)
	case 'S':
		resp = s.search(cmd)
	default:
		log.Printf("Command %02x unrecognised", cmd)
		resp = cmd.EmptyResponse(0xff).
			AppendString("Unsupported")
	}

	if resp == nil {
		resp = cmd.EmptyResponse(0)
	}

	return resp.Write(cmd.out)
}

type SubService struct{}
