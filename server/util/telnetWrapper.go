package util

import (
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/go-telnet"
)

// Wrapper to record idle time & bytes read/written
type TelnetWrapper struct {
	reader telnet.Reader
	writer telnet.Writer
	con    *status.Connection
}

func NewTelnetWrapper(reader telnet.Reader, writer telnet.Writer, con *status.Connection) *TelnetWrapper {
	return &TelnetWrapper{
		reader: reader,
		writer: writer,
		con:    con,
	}
}
func (w *TelnetWrapper) Read(b []byte) (int, error) {
	n, err := w.reader.Read(b)
	if n > 0 {
		w.con.BytesIn(n)
	}
	return n, err
}

func (w *TelnetWrapper) Write(b []byte) (int, error) {
	n, err := w.writer.Write(b)
	if n > 0 {
		w.con.BytesOut(n)
	}
	return n, err
}
