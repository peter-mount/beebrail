package server

import (
	"bufio"
	"errors"
	"io"
)

// Packet represents the actual packet sent over the wire
type Packet struct {
	out     io.Writer // The response writer
	Command byte      // The command code
	Status  byte      // The status when being returned, optional for commands, e.g. sub-command
	Length  uint16    // The length of the payload
	Payload []byte    // The payload
}

func (p *Packet) ErrorPacket(err error) *Packet {
	return p.Response(1, []byte(err.Error())).
		Append(0)
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

func (p *Packet) AppendCString(payload string) *Packet {
	return p.AppendString(payload).Append(0)
}

func (p *Packet) AppendCStrings(cols int, v []string) *Packet {
	for _, h := range v {
		p.AppendCString(h)
	}

	if len(v) < cols {
		p.Append(0)
	}

	return p
}

func (p *Packet) AppendBytes(payload []byte) *Packet {
	p.Payload = append(p.Payload, payload...)
	p.Length = uint16(len(p.Payload))
	return p
}

func (p *Packet) Append(payload ...byte) *Packet {
	return p.AppendBytes(payload)
}

func (p *Packet) PayloadAsString() string {
	return string(p.Payload)
}

// readBlock reads l bytes from in.
// We can't do a simple Read([]byte) call as thats too fast when reading from a BBC over RS423
// so we have to do it with individual reads
func readBlock(in *bufio.Reader, l uint16) ([]byte, error) {
	var buf []byte
	for i := uint16(0); i < l; i++ {
		b, err := in.ReadByte()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b)
	}
	return buf, nil
}

// Read reads a packet from a Reader
func (p *Packet) Read(in *bufio.Reader) error {
	// Read the common header
	b, err := readBlock(in, 4)
	if err != nil {
		return err
	}
	if len(b) != 4 {
		return errors.New("incomplete header")
	}

	p.Command = b[0]
	p.Status = b[1]
	p.Length = uint16(b[2]) | (uint16(b[3]) << 8)

	p.Payload, err = readBlock(in, p.Length)
	if len(p.Payload) != int(p.Length) {
		return errors.New("incomplete payload")
	}

	return nil
}

func (p *Packet) Write(out io.Writer) error {
	b := []byte{p.Command, p.Status, byte(p.Length & 0xff), byte((p.Length >> 8) & 0xff)}
	b = append(b, p.Payload...)
	_, err := out.Write(b)
	return err
}

func (p *Packet) AppendInt16(i int) *Packet {
	return p.Append(
		byte(i&0xff),
		byte((i>>8)&0xff),
	)
}
