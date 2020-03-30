package server

import (
	"errors"
	"fmt"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/go-telnet"
	"github.com/peter-mount/go-telnet/telsh"
	"github.com/peter-mount/golib/kernel"
	"log"
)

type Telnet struct {
	config *Config             // Config file
	server *Server             // the underlying server
	stats  *status.Status      // Status manager
	ports  []*telnetConnection // Active telnet ports
}

type telnetConnection struct {
	parent  *Telnet             // pointer to parent
	config  TelnetConfig        // Config for this instance
	cat     *status.Category    // Status Category for this port
	handler *telsh.ShellHandler // Telnet handler
}

func (t *Telnet) Name() string {
	return "Telnet"
}

func (t *Telnet) Init(k *kernel.Kernel) error {

	service, err := k.AddService(&Config{})
	if err != nil {
		return err
	}
	t.config = (service).(*Config)

	service, err = k.AddService(&Server{})
	if err != nil {
		return err
	}
	t.server = (service).(*Server)

	service, err = k.AddService(&status.Status{})
	if err != nil {
		return err
	}
	t.stats = (service).(*status.Status)

	return nil
}

func (t *Telnet) PostInit() error {
	cat := t.stats.AddCategory("telnet", "Telnet")

	for _, c := range t.config.Telnet {
		if c.Port == 0 {
			return errors.New("invalid telnet port")
		}

		tc := &telnetConnection{
			parent: t,
			config: c,
			cat:    cat,
		}
		tc.cat.Port = c.Port
		t.ports = append(t.ports, tc)
	}

	return nil
}

func (t *Telnet) Start() error {

	for _, c := range t.ports {
		go c.start()
	}

	return nil
}

func (tc *telnetConnection) start() {
	var err error
	port := fmt.Sprintf("%s:%d", tc.config.Interface, tc.config.Port)
	tc.handler = NewShell(tc.parent.server)

	log.Println("Starting telnet", port, "secure", tc.config.Secure)
	if tc.config.Secure {
		err = telnet.ListenAndServeTLS(port, tc.parent.config.Tls.Cert, tc.parent.config.Tls.Key, tc)
	} else {
		err = telnet.ListenAndServe(port, tc)
	}
	if err != nil {
		log.Println(tc.cat.Name, err)
	}
}

func (tc *telnetConnection) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	con := tc.cat.Add(ctx.LocalAddr(), ctx.RemoteAddr())
	defer con.Remove()

	wrapper := &telnetWrapper{
		reader: r,
		writer: w,
		con:    con,
	}

	tc.handler.ServeTELNET(ctx, wrapper, wrapper)
}

// Wrapper to record idle time & bytes read/written
type telnetWrapper struct {
	reader telnet.Reader
	writer telnet.Writer
	con    *status.Connection
}

func (w *telnetWrapper) Read(b []byte) (int, error) {
	n, err := w.reader.Read(b)
	if n > 0 {
		w.con.BytesIn(n)
	}
	return n, err
}

func (w *telnetWrapper) Write(b []byte) (int, error) {
	n, err := w.writer.Write(b)
	if n > 0 {
		w.con.BytesOut(n)
	}
	return n, err
}
