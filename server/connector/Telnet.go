package connector

import (
	"errors"
	"fmt"
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/go-telnet"
	"github.com/peter-mount/go-telnet/telsh"
	"github.com/peter-mount/golib/kernel"
	"github.com/reiver/go-oi"
	"io"
	"log"
	"time"
)

type Telnet struct {
	config *server.Config      // Config file
	server *server.Server      // the underlying server
	stats  *status.Status      // Status manager
	ports  []*telnetConnection // Active telnet ports
}

type telnetConnection struct {
	parent  *Telnet             // pointer to parent
	config  server.TelnetConfig // Config for this instance
	cat     *status.Category    // Status Category for this port
	handler *telsh.ShellHandler // Telnet handler
}

func (t *Telnet) Name() string {
	return "Telnet"
}

func (t *Telnet) Init(k *kernel.Kernel) error {

	service, err := k.AddService(&server.Config{})
	if err != nil {
		return err
	}
	t.config = (service).(*server.Config)

	service, err = k.AddService(&server.Server{})
	if err != nil {
		return err
	}
	t.server = (service).(*server.Server)

	service, err = k.AddService(&status.Status{})
	if err != nil {
		return err
	}
	t.stats = (service).(*status.Status)

	return nil
}

func (t *Telnet) PostInit() error {
	for _, c := range t.config.Telnet {
		if c.Port == 0 {
			return errors.New("invalid telnet port")
		}

		titleFmt := "Telnet Insecure"
		if c.Secure {
			titleFmt = "Telnet Secure"
		}

		tc := &telnetConnection{
			parent: t,
			config: c,
			cat:    t.stats.AddCategory(fmt.Sprintf("telnet%05d", c.Port), titleFmt),
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
	tc.handler = telsh.NewShellHandler()
	tc.handler.Register("test", animateProducer)

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

	tc.handler.ServeTELNET(ctx, w, r)
}

func animateHandler(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {

	for i := 0; i < 20; i++ {
		oi.LongWriteString(stdout, "\r⠋")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠙")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠹")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠸")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠼")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠴")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠦")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠧")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠇")
		time.Sleep(50 * time.Millisecond)

		oi.LongWriteString(stdout, "\r⠏")
		time.Sleep(50 * time.Millisecond)
	}
	oi.LongWriteString(stdout, "\r \r\n")

	return nil
}

func animateProducerFunc(ctx telnet.Context, name string, args ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(animateHandler)
}

var animateProducer = telsh.ProducerFunc(animateProducerFunc)
