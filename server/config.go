package server

import (
	"errors"
	"flag"
	"github.com/jacobsa/go-serial/serial"
	"github.com/peter-mount/golib/kernel"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config holds the yaml config for the server
type Config struct {
	configFile *string // Location of config file

	Serial struct {
		Port                  string            `yaml:"port"`
		BaudRate              uint              `yaml:"baud"`
		DataBits              uint              `yaml:"data"`
		StopBits              uint              `yaml:"stop"`
		Parity                serial.ParityMode `yaml:"parity"`
		FlowControl           bool              `yaml:"flowControl"`
		InterCharacterTimeout uint              `yaml:"interCharacterTimeout"`
		MinimumReadSize       uint              `yaml:"minimumReadSize"`
	} `yaml:"serial"`

	Telnet struct {
		Port int `yaml:"port"` // Telnet port
	} `yaml:"telnet"`

	TelnetSecure struct {
		Port int    `yaml:"port"` // TelnetS port
		Cert string `yaml:"cert"` // Cert.pem path
		Key  string `yaml:"key"`  // Key.pem path
	} `yaml:"telnets"`

	Services struct {
		Reference string `yaml:"reference"` // The ref service url
		LDB       string `yaml:"ldb"`       // The LDB service url
	} `yaml:"services"`
}

func (c *Config) Name() string {
	return "Config"
}

func (c *Config) Init(k *kernel.Kernel) error {
	c.configFile = flag.String("c", "config.yaml", "Configuration file")
	return nil
}

func (c *Config) PostInit() error {
	if c.configFile == nil || *c.configFile == "" {
		return errors.New("config file is required")
	}

	if filename, err := filepath.Abs(*c.configFile); err != nil {
		return err
	} else if in, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else if err := yaml.Unmarshal(in, c); err != nil {
		return err
	}

	if c.Services.LDB == "" {
		return errors.New("ldb service undefined")
	}

	if c.Services.Reference == "" {
		return errors.New("reference service undefined")
	}

	return nil
}
