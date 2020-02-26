package server

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"io"
	"log"
)

type Server struct {
	port      io.ReadWriteCloser // the serial port
	in        *bufio.Reader      // reader on the port
	refClient refclient.DarwinRefClient
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) PostInit() error {
	log.Println("Starting server")

	options := serial.OpenOptions{
		PortName:              "/dev/ttyUSB0",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
		RTSCTSFlowControl:     true,
	}

	port, err := serial.Open(options)
	if err != nil {
		return err
	}
	s.port = port

	s.in = bufio.NewReader(port)

	s.refClient = refclient.DarwinRefClient{Url: "https://ref.prod.a51.li"}
	return nil
}

func (s *Server) Start() error {

	for true {
		_ = s.processCommand()
	}
	return nil
}

func (s *Server) processCommand() error {
	cmd := Packet{}
	err := cmd.Read(s.in)
	if err != nil {
		return err
	}

	var resp *Packet
	switch cmd.Command {
	case 'C':
		resp = s.crs(cmd)
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

	return resp.Write(s.port)
}
