package connector

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/golib/kernel"
	"io"
	"log"
)

// Handles the serial port option
type SerialPort struct {
	config *server.Config      // Config file
	server *server.Server      // the underlying server
	stats  *status.Status      // Status manager
	cat    *status.Category    // Status Category
	ports  []*serialConnection // Active connections
}

type serialConnection struct {
	parent *SerialPort         // pointer to parent
	config server.SerialConfig // Config for this instance
	conn   *status.Connection  // Status connection
	port   io.ReadWriteCloser  // the serial port
	in     *bufio.Reader       // reader on the port
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

	service, err = k.AddService(&status.Status{})
	if err != nil {
		return err
	}
	s.stats = (service).(*status.Status)

	return nil
}

func (s *SerialPort) PostInit() error {
	s.cat = s.stats.AddCategory("serial", "Direct Serial")

	for _, c := range s.config.Serial {
		if c.BaudRate == 0 {
			c.BaudRate = 9600
		}

		if c.DataBits == 0 {
			c.DataBits = 8
			c.StopBits = 1
		}

		s.ports = append(s.ports, &serialConnection{
			parent: s,
			config: c,
		})
	}

	return nil
}

func (s *SerialPort) Start() error {
	for _, c := range s.ports {
		err := c.start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serialConnection) start() error {
	port, err := serial.Open(serial.OpenOptions{
		PortName:              s.config.Port,
		BaudRate:              s.config.BaudRate,
		DataBits:              s.config.DataBits,
		StopBits:              s.config.StopBits,
		ParityMode:            s.config.Parity,
		RTSCTSFlowControl:     s.config.FlowControl,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	})
	if err != nil {
		return err
	}

	s.port = port
	s.in = bufio.NewReader(port)

	go s.run()

	return nil
}

func (s *serialConnection) run() {
	// On exit close the port
	defer s.port.Close()

	// Register the connection in stats
	s.conn = s.parent.cat.Add(nil, nil)
	s.conn.Name = s.config.Port

	defer s.conn.Remove()

	for true {
		packet, err := s.parent.server.ReadPacket(s.in, s.port)

		if err == nil {
			err = s.parent.server.ProcessPacket(packet)
		}

		if err != nil && err != io.EOF {
			log.Println("Serial", err)
		}
	}
}
