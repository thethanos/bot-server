package config

import (
	"github.com/pelletier/go-toml"
)

type Mode string

const (
	DEBUG   Mode = "debug"
	RELEASE Mode = "release"
)

type Config struct {
	Mode        Mode
	ReleasePort int64
	DebugPort   int64
	ImagePrefix string
	PsqlHost    string
	PsqlPort    int64
	PsqlUser    string
	PsqlPass    string
	PsqlDb      string
}

func Load(path string) (*Config, error) {

	cfg, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		Mode:        Mode(cfg.Get("bot-server.mode").(string)),
		ReleasePort: cfg.Get("bot-server.release_port").(int64),
		DebugPort:   cfg.Get("bot-server.debug_port").(int64),
		ImagePrefix: cfg.Get("bot-server.image_prefix").(string),
		PsqlHost:    cfg.Get("postgres.host").(string),
		PsqlPort:    cfg.Get("postgres.port").(int64),
		PsqlUser:    cfg.Get("postgres.user").(string),
		PsqlPass:    cfg.Get("postgres.password").(string),
		PsqlDb:      cfg.Get("postgres.dbname").(string),
	}, nil
}
