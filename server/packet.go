package server

import (
	"bufio"
	"io"
)

type Packet struct {
	Command byte   // The command code
	Status  byte   // The status when being returned
	Length  uint16 // The length of the payload
	Payload []byte // The payload
}

func (p *Packet) ErrorPacket(err error) *Packet {
	return p.Response(1, []byte(err.Error()))
}

func (p *Packet) EmptyResponse(status byte) *Packet {
	return p.Response(status, []byte{})
}

func (p *Packet) Response(status byte, payload []byte) *Packet {
	return &Packet{
		Command: p.Command,
		Status:  status,
		Length:  uint16(len(payload)),
		Payload: payload,
	}
}

func (p *Packet) AppendString(payload string) *Packet {
	return p.Append([]byte(payload)...)
}

func (p *Packet) Append(payload ...byte) *Packet {
	p.Payload = append(p.Payload, payload...)
	p.Length = uint16(len(p.Payload))
	return p
}

func (p *Packet) PayloadAsString() string {
	return string(p.Payload)
}

// Read reads a packet from a Reader
func (p *Packet) Read(in *bufio.Reader) error {
	var err error
	p.Command, err = in.ReadByte()
	if err != nil {
		return err
	}

	// Status byte, ignore for command invocation
	p.Status, err = in.ReadByte()
	if err != nil {
		return err
	}

	// Read in packet length
	l1, err := in.ReadByte()
	if err != nil {
		return err
	}

	l2, err := in.ReadByte()
	if err != nil {
		return err
	}

	p.Length = uint16(l1) | (uint16(l2) << 8)

	// Note we can't do the following as the BBC isn't fast enough!
	// l, err := s.in.Read(packet)
	for i := uint16(0); i < p.Length; i++ {
		b, err := in.ReadByte()
		if err != nil {
			return err
		}
		p.Payload = append(p.Payload, b)
	}

	return nil
}

func (p *Packet) Write(out io.Writer) error {
	b := []byte{p.Command, p.Status, byte(p.Length & 0xff), byte((p.Length >> 8) & 0xff)}
	_, err := out.Write(b)
	if err != nil {
		return err
	}
	_, err = out.Write(p.Payload)
	return err
}
