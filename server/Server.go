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
		cmd, err := s.in.ReadString(0)
		if err == nil && len(cmd) > 0 {
			param := cmd[1 : len(cmd)-1]
			log.Println(cmd[0:1], param)

			switch cmd[0] {
			//case 'C':
			//	s.crs(s[1..])
			case 'S':
				s.search(param)
			default:
				log.Println("Command unrecognised")
			}
		}
	}
	return nil
}
