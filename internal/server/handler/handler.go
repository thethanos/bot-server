package handler

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	"bot/internal/minioadapter"
)

type Handler struct {
	logger       logger.Logger
	cfg          *config.Config
	DBAdapter    *dbadapter.DBAdapter
	MinIOAdapter *minioadapter.MinIOAdapter
}

func NewHandler(logger logger.Logger, cfg *config.Config, DBAdapter *dbadapter.DBAdapter, MinIOAdapter *minioadapter.MinIOAdapter) *Handler {
	return &Handler{
		logger:       logger,
		cfg:          cfg,
		DBAdapter:    DBAdapter,
		MinIOAdapter: MinIOAdapter,
	}
}
