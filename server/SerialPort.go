package server

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
)

// Handles the serial port option
type SerialPort struct {
	server *Server            // the underlying server
	port   io.ReadWriteCloser // the serial port
	in     *bufio.Reader      // reader on the port
}

func (s *Server) initSerialPort() *SerialPort {
	port := &SerialPort{server: s}

	c := s.config

	if c.Serial.BaudRate == 0 {
		c.Serial.BaudRate = 9600
	}

	if c.Serial.DataBits == 0 {
		c.Serial.DataBits = 8
		c.Serial.StopBits = 1
		c.Serial.MinimumReadSize = 0
	}

	if c.Serial.InterCharacterTimeout == 0 {
		c.Serial.InterCharacterTimeout = 100
	}

	return port
}

func (s *SerialPort) Start() error {
	log.Println("Starting serialPort")

	c := s.server.config.Serial
	port, err := serial.Open(serial.OpenOptions{
		PortName:              c.Port,
		BaudRate:              c.BaudRate,
		DataBits:              c.DataBits,
		StopBits:              c.StopBits,
		ParityMode:            c.Parity,
		RTSCTSFlowControl:     c.FlowControl,
		InterCharacterTimeout: c.InterCharacterTimeout,
		MinimumReadSize:       c.MinimumReadSize,
	})
	if err != nil {
		return err
	}
	s.port = port

	s.in = bufio.NewReader(port)

	return nil
}

func (s *SerialPort) Run() error {
	for true {
		packet, err := s.server.ReadPacket(s.in, s.port)

		if err == nil {
			err = s.server.ProcessPacket(packet)
		}

		if err != nil && err != io.EOF {
			log.Println("Serial", err)
		}
	}

	return nil
}
