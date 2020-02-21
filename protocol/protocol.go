package protocol

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
)

type Protocol struct {
	port io.ReadWriteCloser

	in *bufio.Reader
}

func (p *Protocol) Name() string {
	return "Protocol"
}

func (p *Protocol) PostInit() error {
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
	p.port = port

	p.in = bufio.NewReader(port)

	return nil
}

func (p *Protocol) Start() error {

	p.run()
	return nil
}

func (p *Protocol) run() {
	for true {
		s, err := p.in.ReadString(0)
		if err != nil {
			//log.Println(err)
		} else {
			p.exec(s)
		}
	}
}

func (p *Protocol) exec(s string) {
	log.Println(s)
}
