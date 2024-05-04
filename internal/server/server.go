package server

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	"bot/internal/minioadapter"
	"bot/internal/server/handler"
	corsMiddleware "bot/internal/server/middleware"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func NewServer(logger logger.Logger, cfg *config.Config, DBAdapter *dbadapter.DBAdapter, MinIOAdapter *minioadapter.MinIOAdapter) (*http.Server, error) {

	handler := handler.NewHandler(logger, cfg, DBAdapter, MinIOAdapter)
	docHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "swagger.yaml"}, nil)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/cities", handler.GetCities)
	getRouter.HandleFunc("/services/categories", handler.GetServiceCategories)
	getRouter.HandleFunc("/services", handler.GetServices)
	getRouter.HandleFunc("/masters/bot", handler.GetMastersBot)
	getRouter.HandleFunc("/masters/admin", handler.GetMastersAdmin)
	getRouter.HandleFunc("/masters/{master_id}", handler.GetMaster)
	getRouter.HandleFunc("/masters/{master_id}/images", handler.GetMasterImages)
	getRouter.Handle("/docs", docHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("/bot-server/docs")))

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/cities", handler.SaveCity)
	postRouter.HandleFunc("/services/categories", handler.SaveServiceCategory)
	postRouter.HandleFunc("/services", handler.SaveService)
	postRouter.HandleFunc("/masters", handler.SaveMaster)
	postRouter.HandleFunc("/masters/{master_id}/images", handler.SaveMasterImage)
	postRouter.HandleFunc("/masters/approve/{master_id}", handler.ApproveMaster)

	putHandler := router.Methods(http.MethodPut).Subrouter()
	putHandler.HandleFunc("/cities", handler.UpdateCity)
	putHandler.HandleFunc("/services/categories", handler.UpdateServCategory)
	putHandler.HandleFunc("/services", handler.UpdateService)
	putHandler.HandleFunc("/masters", handler.UpdateMaster)
	putHandler.HandleFunc("/masters/{master_id}/images/{image_name}", handler.UpdateMasterImage)

	deleteHandler := router.Methods(http.MethodDelete).Subrouter()
	deleteHandler.HandleFunc("/cities/{city_id}", handler.DeleteCity)
	deleteHandler.HandleFunc("/services/categories/{category_id}", handler.DeleteServCategory)
	deleteHandler.HandleFunc("/services/{service_id}", handler.DeleteService)
	deleteHandler.HandleFunc("/masters/{master_id}", handler.DeleteMaster)
	deleteHandler.HandleFunc("/masters/{master_id}/images/{image_name}", handler.DeleteMasterImage)

	return &http.Server{
		Handler: corsMiddleware.CorsMiddlware(router),
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}, nil
}
