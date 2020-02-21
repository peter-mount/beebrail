package server

import (
	"github.com/jacobsa/go-serial/serial"
	"github.com/peter-mount/golib/kernel"
	"io"
	"log"
)

type Server struct {
	port io.ReadWriteCloser
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) Init(k *kernel.Kernel) error {
	return nil
}

func (s *Server) Start() error {
	log.Println("Starting server")
	options := serial.OpenOptions{
		PortName:              "/dev/ttyUSB0",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		return err
	}
	s.port = port

	s.port.Write([]byte("Hello\r\n"))

	return nil
}
