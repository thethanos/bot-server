package handler

import (
	"bot/internal/entities"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Update city
// @Description Change the city name
// @Tags City
// @Param city body entities.City true "City id and name"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {string} string "Error message"
// @Failure 404 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /cities [put]
func (h *Handler) UpdateCity(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::UpdateCity::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	city := &entities.City{}
	if err := json.Unmarshal(body, city); err != nil {
		h.logger.Error("server::UpdateCity::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DBAdapter.UpdateCity(city); err != nil {
		h.logger.Error("server::UpdateCity::UpdateCity")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	h.logger.Info("Response sent")
}

// @Summary Update service category
// @Description Change the service categiry name
// @Tags Service
// @Param service body entities.ServiceCategory true "Service category id and name"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {string} string "Error message"
// @Failure 404 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services/categories [put]
func (h *Handler) UpdateServCategory(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::UpdateServCategory::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	category := &entities.ServiceCategory{}
	if err := json.Unmarshal(body, category); err != nil {
		h.logger.Error("server::UpdateServCategory::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DBAdapter.UpdateServCategory(category); err != nil {
		h.logger.Error("server::UpdateServCategory::UpdateServCategory")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	h.logger.Info("Response sent")
}

// @Summary Update service
// @Description Change the service name or category
// @Tags Service
// @Param service body entities.Service true "Service category id, category name, id and name"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {string} string "Error message"
// @Failure 404 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services [put]
func (h *Handler) UpdateService(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::UpdateService::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	service := &entities.Service{}
	if err := json.Unmarshal(body, service); err != nil {
		h.logger.Error("server::UpdateService::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DBAdapter.UpdateService(service); err != nil {
		h.logger.Error("server::UpdateService::UpdateService")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	h.logger.Info("Response sent")
}

// @Summary Update master
// @Description Update master data in the system
// @Tags Master
// @Param service body entities.MasterLong true "Master data"
// @Accept json
// @Produce json
// @Success 200 {object} ID "ID of the updated master"
// @Failure 400 {string} string "Error message"
// @Failure 404 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [put]
func (h *Handler) UpdateMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::UpdateMaster::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	master := &entities.MasterLong{}
	if err := json.Unmarshal(body, master); err != nil {
		h.logger.Error("server::UpdateMaster::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	validator := validator.New()
	if err := validator.Struct(master); err != nil {
		h.logger.Error("server::UpdateMaster::Struct", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DBAdapter.UpdateMaster(master); err != nil {
		h.logger.Error("server::UpdateMaster::UpdateMaster")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, master.ID))); err != nil {
		h.logger.Errorf("server::ApproveMaster::Write: %s", err.Error())
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Update master image
// @Description Update an image of a master in the system
// @Tags Master
// @Param master_id path string true "ID of the master"
// @Param image_name path string true "Name of the image"
// @Param file formData file true "Updated image"
// @Accept multipart/form-data
// @Produce json
// @Success 204
// @Failure 400 {string} string "Error message"
// @Failure 404 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id}/images/{image_name} [put]
func (h *Handler) UpdateMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]
	imageName := params["image_name"]

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("server::UpdateMasterImage::ParseMultipartForm", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	formFile, meta, err := req.FormFile("file")
	if err != nil {
		h.logger.Error("server::UpdateMasterImage::FormFile", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	if err := h.MinIOAdapter.PutMasterImage(masterID, uuid.NewString(), formFile, meta.Size, meta.Header.Get("Content-Type")); err != nil {
		h.logger.Error("server::UpdateMasterImage::PutMasterImage", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.MinIOAdapter.DeleteMasterImage(masterID, imageName); err != nil {
		h.logger.Error("server::UpdateMasterImage::DeleteMasterImage", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	h.logger.Info("Response sent")
}
