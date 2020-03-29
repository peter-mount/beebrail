package server

import (
	"bufio"
	"errors"
	"github.com/peter-mount/golib/kernel"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	ldbclient "github.com/peter-mount/nre-feeds/ldb/client"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
)

type Server struct {
	config     *Config                   // Config file
	refClient  refclient.DarwinRefClient // ref api
	ldbClient  ldbclient.DarwinLDBClient // ldb api
	handlers   []ServerHandler           // slice of handlers
	serialPort *SerialPort               // serialPort handler
}

// ServerHandler interface implemented by handlers, e.g. direct Serial port, Telnet etc
type ServerHandler interface {
	Start() error
	Run() error
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

	log.Println(s.config)
	log.Println(s.config.Serial)
	log.Println(s.config.Serial.Port)
	if s.config.Serial.Port != "" {
		s.handlers = append(s.handlers, s.initSerialPort())
	}

	return nil
}

func (s *Server) Start() error {

	// Can't start if we have no handlers
	if len(s.handlers) == 0 {
		return errors.New("No interfaces were defined")
	}

	for _, h := range s.handlers {
		if err := h.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Run() error {
	// Run a Group in JustErrors mode, terminate on the first handler to exit with no error
	var group errgroup.Group

	for _, h := range s.handlers {
		group.Go(h.Run)
	}

	return group.Wait()
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
