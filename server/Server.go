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
	cmd, err := s.in.ReadByte()
	if err != nil {
		return err
	}
	log.Printf("%02x", cmd)

	// Status byte, ignore for command invocation
	st, err := s.in.ReadByte()
	if err != nil {
		return err
	}
	log.Printf("%02x", st)

	// Read in packet length
	b, err := s.in.ReadByte()
	if err != nil {
		return err
	}
	log.Printf("%02x", b)
	l := int(b)

	b, err = s.in.ReadByte()
	if err != nil {
		return err
	}
	log.Printf("%02x", b)
	l = l | (int(b) << 8)

	var packet []byte
	// Note we can't do the following as the BBC isn't fast enough!
	// l, err := s.in.Read(packet)
	for i := 0; i < l; i++ {
		b, err = s.in.ReadByte()
		if err != nil {
			return err
		}
		packet = append(packet, b)
	}

	log.Printf("%02x Received %04x expected %04x", cmd, len(packet), l)

	switch cmd {
	case 'S':
		s.search(string(packet))
	default:
		log.Printf("Command %02x unrecognised", cmd)
	}

	return nil
}
