package server

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/entities"
	"bot/internal/logger"
	"bot/internal/minioadapter"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

// @Summary Get cities
// @Description Get all available cities
// @Tags City
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Accept json
// @Produce json
// @Success 200 {array} entities.City
// @Failure 500 {string} string "Error message"
// @Router /cities [get]
func (h *Handler) GetCities(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cities, err := h.DBAdapter.GetCities("", page, limit)
	if err != nil {
		h.logger.Error("server::GetCities::GetCities", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	cityList, err := json.Marshal(&cities)
	if err != nil {
		h.logger.Error("server::GetCities::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(cityList); err != nil {
		h.logger.Error("server::GetCities::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get service categories
// @Description Get all available service categories
// @Tags Service
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Acept json
// @Produce json
// @Success 200 {array} entities.ServiceCategory
// @Failure 500 {string} string "Error message"
// @Router /services/categories [get]
func (h *Handler) GetServiceCategories(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := h.DBAdapter.GetServCategories("", page, limit)
	if err != nil {
		h.logger.Error("server::GetCategories::GetCategories", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryList, err := json.Marshal(&categories)
	if err != nil {
		h.logger.Error("server::GetCategories::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(categoryList); err != nil {
		h.logger.Error("server::GetServiceCategories::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get services
// @Description Get all available services, filters by category_id if provided
// @Tags Service
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Param category_id query string false "ID of the service category"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Service
// @Failure 500 {string} string "Error message"
// @Router /services [get]
func (h *Handler) GetServices(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	services, err := h.DBAdapter.GetServices(query.Get("category_id"), "", page, limit)
	if err != nil {
		h.logger.Error("server::GetServices::GetServices", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceList, err := json.Marshal(&services)
	if err != nil {
		h.logger.Error("server::GetServices::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(serviceList); err != nil {
		h.logger.Error("server::GetServices::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get masters
// @Description Get all available masters for the selected city and the service
// @Tags Master
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Param city_id query string false "ID of the selected city"
// @Param service_id query string false "ID of the seleted service"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Master
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [get]
func (h *Handler) GetMasters(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetMasters::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetMasters::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	masters, err := h.DBAdapter.GetMasters(query.Get("city_id"), "", query.Get("service_id"), page, limit)
	if err != nil {
		h.logger.Error("server::GetMasters::GetMasters", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mastersResp, err := json.Marshal(masters)
	if err != nil {
		h.logger.Error("server::GetMasters::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(mastersResp); err != nil {
		h.logger.Error("server::GetMasters::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get master
// @Description Get the master by the given ID
// @Tags Master
// @Param master_id path string true "ID of the master"
// @Accept json
// @Produce json
// @Success 200 {object} entities.MasterRegForm
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id} [get]
func (h *Handler) GetMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	master, err := h.DBAdapter.GetMaster(masterID)
	if err != nil {
		h.logger.Error("server::GetMaster::GetMaster")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	masterResp, err := json.Marshal(master)
	if err != nil {
		h.logger.Error("server::GetMaster::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(masterResp); err != nil {
		h.logger.Error("server::GetMaster::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

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
	h.logger.Infof("Request received: %s", req.URL)

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
	h.logger.Infof("Request received: %s", req.URL)

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
	h.logger.Infof("Request received: %s", req.URL)

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

// @Summary Save master registration form
// @Description Save registration form of a new master
// @Tags Master
// @Param form body entities.MasterRegForm true "Registration form"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [post]
func (h *Handler) SaveMasterRegForm(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	regForm := &entities.MasterRegForm{}
	if err := json.Unmarshal(body, regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	validator := validator.New()
	if err := validator.Struct(regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Struct", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.DBAdapter.SaveMasterRegForm(regForm)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::SaveMasterRegForm", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.MinIOAdapter.MakeBucket(id); err != nil {
		h.logger.Error("server::SaveMasterRegForm::MakeBucket", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// temporary, while the approvement mechanism is not integrated
	if _, err := h.DBAdapter.SaveMaster(id); err != nil {
		h.logger.Error("server::SaveMasterRegForm::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Write", err)
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
// @Router /masters/images/{master_id} [post]
func (h *Handler) SaveMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

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

	if err := h.MinIOAdapter.PutObject(masterID, meta.Filename, formFile, meta.Size, meta.Header.Get("Content-Type")); err != nil {
		h.logger.Error("server::SaveMasterImage::PutObject", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DBAdapter.SaveMasterImage(masterID, meta.Filename); err != nil {
		h.logger.Error("server::SaveMasterImage::SaveMasterImage", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "url" : "%s" }`, meta.Filename))); err != nil {
		h.logger.Error("server::SaveMasterImage::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Approve master
// @Description Approve and save master in the system
// @Tags Master
// @Param master_id path string true "ID of the approved master"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the approved master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/approve/{maser_id} [post]
func (h *Handler) ApproveMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	if _, err := h.DBAdapter.SaveMaster(masterID); err != nil {
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
	h.logger.Infof("Request received: %s", req.URL)

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
	h.logger.Infof("Request received: %s", req.URL)

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
	h.logger.Infof("Request received: %s", req.URL)

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

// @Summary Delete city
// @Description Delete a city from the system
// @Tags City
// @Param city_id path string true "ID of the city"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /cities/{city_id} [delete]
func (h *Handler) DeleteCity(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)

	if err := h.DBAdapter.DeleteCity(params["city_id"]); err != nil {
		h.logger.Error("server::DeleteCity::DeleteCity", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}

// @Summary Delete service category
// @Description Delete a service category along with all its services from the system
// @Tags Service
// @Param category_id path string true "ID of the service category"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services/categories/{category_id} [delete]
func (h *Handler) DeleteServCategory(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)

	if err := h.DBAdapter.DeleteServCategory(params["category_id"]); err != nil {
		h.logger.Error("server::DeleteServCategory::DeleteServCategory", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}

// @Summary Delete service
// @Description Delete a service from the system
// @Tags Service
// @Param service_id path string true "ID of the service"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services/{service_id} [delete]
func (h *Handler) DeleteService(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)

	if err := h.DBAdapter.DeleteService(params["service_id"]); err != nil {
		h.logger.Error("server::DeleteService::DeleteService", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}

// @Summary Delete master
// @Description Delete a master from the system
// @Tags Master
// @Param master_id path string true "ID of the master"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id} [delete]
func (h *Handler) DeleteMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	if err := h.DBAdapter.DeleteMaster(masterID); err != nil {
		h.logger.Error("server::DeleteMaster::DeleteMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.MinIOAdapter.DeleteBucket(masterID); err != nil {
		h.logger.Error("server::DeleteMaster::DeleteBucket", err)
	}

	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}
