package server

import (
	"errors"
	"flag"
	"github.com/peter-mount/golib/kernel"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config holds the yaml config for the server
type Config struct {
	configFile *string // Location of config file

	Telnet []TelnetConfig `yaml:"telnet"` // Telnet connections

	Tls struct {
		Cert string `yaml:"cert"` // Cert.pem path
		Key  string `yaml:"key"`  // Key.pem path
	}

	Services struct {
		Reference string `yaml:"reference"` // The ref service url
		LDB       string `yaml:"ldb"`       // The LDB service url
	} `yaml:"services"`
}

type TelnetConfig struct {
	Interface string `yaml:"interface"` // Interface "" for any
	Port      uint16 `yaml:"port"`      // Port
	Secure    bool   `yaml:"secure"`    // Secure or insecure
	API       bool   `yaml:"api"`       // True for shell intended for computer rather than human
	Shell     struct {
		Prompt         string `yaml:"prompt"`         // Command prompt
		WelcomeMessage string `yaml:"welcomeMessage"` // Welcome message
		ExitMessage    string `yaml:"exitMessage"`    // Exit message
	} `yaml:"shell"`
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
