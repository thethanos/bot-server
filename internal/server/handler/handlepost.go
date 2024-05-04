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

// @Summary Save city
// @Description Save a new city in the system
// @Tags City
// @Param name body Name true "City name"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new city"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /cities [post]
func (h *Handler) SaveCity(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveCity::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	city := &entities.City{}
	if err := json.Unmarshal(body, city); err != nil {
		h.logger.Error("server::SaveCity::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.DBAdapter.SaveCity(city.Name)
	if err != nil {
		h.logger.Error("server::SaveCity::SaveCity", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveCity::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save service category
// @Description Save a new service category in the system
// @Tags Service
// @Param name body Name true "Service category name"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new service category"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services/categories [post]
func (h *Handler) SaveServiceCategory(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveServiceCategory::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	serviceCategory := &entities.ServiceCategory{}
	if err := json.Unmarshal(body, serviceCategory); err != nil {
		h.logger.Error("server::SaveServiceCategory::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.DBAdapter.SaveServiceCategory(serviceCategory.Name)
	if err != nil {
		h.logger.Error("server::SaveServiceCategory::SaveServiceCategory", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveServiceCategory::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save service
// @Description Save a new service in the system
// @Tags Service
// @Param service body entities.Service true "New service"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new service"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services [post]
func (h *Handler) SaveService(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveService::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	service := &entities.Service{}
	if err := json.Unmarshal(body, service); err != nil {
		h.logger.Error("server::SaveService::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.DBAdapter.SaveService(service.Name, service.CatID)
	if err != nil {
		h.logger.Error("server::SaveService::SaveService", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveService::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save master
// @Description Save new master in the system
// @Tags Master
// @Param form body entities.Master true "Master data"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [post]
func (h *Handler) SaveMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveMaster::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	master := &entities.Master{}
	if err := json.Unmarshal(body, master); err != nil {
		h.logger.Error("server::SaveMaster::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	validator := validator.New()
	if err := validator.Struct(master); err != nil {
		h.logger.Error("server::SaveMaster::Struct", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.DBAdapter.SaveMaster(master)
	if err != nil {
		h.logger.Error("server::SaveMaster::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.MinIOAdapter.MakeBucket(id); err != nil {
		h.logger.Error("server::SaveMaster::MakeBucket", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// temporary, while the approvement mechanism is not integrated
	if err := h.DBAdapter.ApproveMaster(id); err != nil {
		h.logger.Error("server::SaveMaster::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveMaster::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save master's image
// @Description Save the image that was attached to the registration form
// @Tags Master
// @Param master_id path string true "ID of a master, whose picture is uploaded"
// @Param file formData file true "Image to upload"
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} URL "URL of the saved picture"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id}/images [post]
func (h *Handler) SaveMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("server::SaveMasterImage::ParseMultipartForm", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	formFile, meta, err := req.FormFile("file")
	if err != nil {
		h.logger.Error("server::SaveMasterImage::FormFile", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	newImageName := uuid.NewString()
	if err := h.MinIOAdapter.PutMasterImage(masterID, newImageName, formFile, meta.Size, meta.Header.Get("Content-Type")); err != nil {
		h.logger.Error("server::SaveMasterImage::PutMasterImage", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "url" : "%s" }`, newImageName))); err != nil {
		h.logger.Error("server::SaveMasterImage::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Approve master
// @Description Approve master to be listed in the system
// @Tags Master
// @Param master_id path string true "ID of the approved master"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the approved master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/approve/{maser_id} [post]
func (h *Handler) ApproveMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	if err := h.DBAdapter.ApproveMaster(masterID); err != nil {
		h.logger.Error("server::ApproveMaster::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, masterID))); err != nil {
		h.logger.Error("server::ApproveMaster::Write", err)
		return
	}
	h.logger.Info("Response sent")
}
