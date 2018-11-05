package config

import (
	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
)

type (
	Config struct {
		Debug            bool
		ConfigFile       string
		IDImageUploadDir string
		*Srv
		*DB
		*IDCheckAPI
		*SMSAPI
	}

	Srv struct {
		Host string
		Port string
	}

	DB struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}

	IDCheckAPI struct {
		Url      string
		Username string
		Password string
	}

	SMSAPI struct {
		URL      string
		Username string
		Password string
		PID      string
	}
)

// LoadFromIni load config from ini override default config
func (config *Config) LoadFromIni() (err error) {
	return ini.MapTo(config, config.ConfigFile)
}

// Load load config from command line param
func (config *Config) Load(c *cli.Context) (err error) {

	if c.String("config") != "" {
		config.ConfigFile = c.String("config")
		if err = config.LoadFromIni(); err != nil {
			return
		}
	}

	if c.Bool("debug") {
		config.Debug = true
	}

	if port := c.String("port"); "" != port {
		config.Srv.Port = port
	}

	if bind := c.String("bind"); "" != bind {
		config.Srv.Host = bind
	}

	return
}
