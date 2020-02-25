package server

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
)

type Server struct {
	port io.ReadWriteCloser

	in *bufio.Reader
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

	return nil
}

func (s *Server) Start() error {

	for true {
		_ = s.processCommand()
	}
	return nil
}

func (s *Server) processCommand() error {
	cmd := &Packet{}
	err := cmd.Read(s.in)
	if err != nil {
		return err
	}

	var resp *Packet
	switch cmd.Command {
	case 'S':
		resp = s.search(cmd)
	default:
		log.Printf("Command %02x unrecognised", cmd)
	}

	if resp == nil {
		resp = cmd.EmptyResponse(0)
	}

	return resp.Write(s.port)
}
