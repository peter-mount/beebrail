package server

import (
	"errors"
	"fmt"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/beebrail/server/util"
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
	tc.handler = tc.NewShell(tc.parent.server)

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

const (
	KEY_CONTEXT = "context"
)

func (tc *telnetConnection) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	con := tc.cat.Add(ctx.LocalAddr(), ctx.RemoteAddr())
	defer con.Remove()

	// Default settings
	m := ctx.UserData()
	shellCtx := &ShellContext{
		connection: con,
		tableStyle: util.Plain,
		userData:   m,
		api:        tc.config.API,
	}
	(*m)[KEY_CONTEXT] = shellCtx

	// The appropriate ResponseWriter
	if tc.config.API {
		shellCtx.responseWriter = util.NewAPIResponseWriter(w)
	} else {
		shellCtx.responseWriter = util.NewHumanResponseWriter(w)
	}

	wrapper := util.NewTelnetWrapper(r, w, con)

	tc.handler.ServeTELNET(ctx, wrapper, wrapper)
}
