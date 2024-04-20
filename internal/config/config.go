package config

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	Port        int64
	ImagePrefix string
	PsqlHost    string
	PsqlPort    int64
	PsqlUser    string
	PsqlPass    string
	PsqlDb      string
	MinIOHost   string
	MinIOPort   int64
	MinIOUser   string
	MinIOPass   string
}

func Load(path string) (*Config, error) {

	cfg, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        cfg.Get("bot-server.port").(int64),
		ImagePrefix: cfg.Get("bot-server.image_prefix").(string),
		PsqlHost:    cfg.Get("postgres.host").(string),
		PsqlPort:    cfg.Get("postgres.port").(int64),
		PsqlUser:    cfg.Get("postgres.user").(string),
		PsqlPass:    cfg.Get("postgres.password").(string),
		PsqlDb:      cfg.Get("postgres.dbname").(string),
		MinIOHost:   cfg.Get("minio.host").(string),
		MinIOPort:   cfg.Get("minio.port").(int64),
		MinIOUser:   cfg.Get("minio.user").(string),
		MinIOPass:   cfg.Get("minio.password").(string),
	}, nil
}
