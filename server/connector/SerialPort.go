package connector

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/go-telnet"
	"github.com/peter-mount/golib/kernel"
	"io"
	"log"
	"time"
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

	// Dependency is only to start after telnet
	_, err = k.AddService(&Telnet{})

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
		PortName:              s.config.Device,
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
	//s.in = bufio.NewReader(port)

	go s.run()

	return nil
}

func (s *serialConnection) run() {
	// On exit close the port
	defer s.port.Close()

	err := telnet.DialToAndCall(s.config.Port, s)
	if err != nil {
		log.Println("Serial", err)
	}
}

func (s *serialConnection) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	conn := s.parent.cat.Add(ctx.LocalAddr(), ctx.RemoteAddr())
	conn.Name = s.config.Device
	defer conn.Remove()

	_, _ = s.port.Write([]byte{0xFF, 0xFB, 0x01}) // IAC WILL ECHO

	go pipe(s.port, r)
	pipe(w, s.port)

	// Wait a bit to receive data from the server (that we would send to io.Stdout).
	time.Sleep(3 * time.Millisecond)
}

func pipe(w io.Writer, r io.Reader) {

	// Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	var buffer [1]byte
	p := buffer[:]

	for {
		// Read 1 byte.
		n, err := r.Read(p)
		if err != nil && err != io.EOF {
			break
		} else if n > 0 {
			_, _ = w.Write(p)
		} else {
			// Wait a short while to allow some input to arrive
			time.Sleep(3 * time.Millisecond)
		}
	}

}
