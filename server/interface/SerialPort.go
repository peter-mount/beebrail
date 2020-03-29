package _interface

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/golib/kernel"
	"io"
	"log"
)

// Handles the serial port option
type SerialPort struct {
	config *server.Config     // Config file
	server *server.Server     // the underlying server
	port   io.ReadWriteCloser // the serial port
	in     *bufio.Reader      // reader on the port
}

func (s *SerialPort) Name() string {
	return "SerialPort"
}

func (s *SerialPort) Init(k *kernel.Kernel) error {

	service, err := k.AddService(&server.Config{})
	if err != nil {
		return err
	}
	s.config = (service).(*server.Config)

	service, err = k.AddService(&server.Server{})
	if err != nil {
		return err
	}
	s.server = (service).(*server.Server)

	return nil
}

func (s *SerialPort) PostInit() error {
	c := s.config.Serial

	if c.BaudRate == 0 {
		c.BaudRate = 9600
	}

	if c.DataBits == 0 {
		c.DataBits = 8
		c.StopBits = 1
		c.MinimumReadSize = 0
	}

	if c.InterCharacterTimeout == 0 {
		c.InterCharacterTimeout = 100
	}

	return nil
}

func (s *SerialPort) Start() error {
	log.Println("Starting serialPort")

	c := s.config.Serial
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

	// Start the main thread in a separate goroutine
	go s.mainThread()

	return nil
}

func (s *SerialPort) mainThread() {
	for true {
		packet, err := s.server.ReadPacket(s.in, s.port)

		if err == nil {
			err = s.server.ProcessPacket(packet)
		}

		if err != nil && err != io.EOF {
			log.Println("Serial", err)
		}
	}
}
