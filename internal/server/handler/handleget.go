package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

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
// @Description Get all available masters for the selected city and the service. Used by the bot.
// @Tags Master
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Param city_id query string false "ID of the selected city"
// @Param service_id query string false "ID of the seleted service"
// @Accept json
// @Produce json
// @Success 200 {array} entities.MasterShort
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/bot [get]
func (h *Handler) GetMastersBot(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetMastersBot::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetMastersBot::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	masters, err := h.DBAdapter.GetMastersBot(query.Get("city_id"), "", query.Get("service_id"), page, limit)
	if err != nil {
		h.logger.Error("server::GetMastersBot::GetMastersBot", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for index, _ := range masters {
		masters[index].Images = h.MinIOAdapter.GetMasterImagesURLs(masters[index].ID)
	}

	mastersResp, err := json.Marshal(masters)
	if err != nil {
		h.logger.Error("server::GetMastersBot::Marshal", err)
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

// @Summary Get masters
// @Description Get all available masters. Used by control panel.
// @Tags Master
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of items for pagination"
// @Accept json
// @Produce json
// @Success 200 {array} entities.MasterShort
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/admin [get]
func (h *Handler) GetMastersAdmin(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	query := req.URL.Query()
	page, err := getParam[int](query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetMastersAdmin::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParam[int](query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetMastersAdmin::getParam[int]", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	masters, err := h.DBAdapter.GetMastersAdmin(page, limit)
	if err != nil {
		h.logger.Error("server::GetMastersAdmin::GetMastersAdmin", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mastersResp, err := json.Marshal(masters)
	if err != nil {
		h.logger.Error("server::GetMastersAdmin::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(mastersResp); err != nil {
		h.logger.Error("server::GetMastersAdmin::Write", err)
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
// @Success 200 {object} entities.MasterLong
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id} [get]
func (h *Handler) GetMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

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

// @Summary Get master images
// @Description Gat all the images provided by master
// @Tags Master
// @Param master_id path string true "ID of the master"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Image
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id}/images [get]
func (h *Handler) GetMasterImages(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	images := h.MinIOAdapter.GetMasterImages(masterID)

	imagesResp, err := json.Marshal(images)
	if err != nil {
		h.logger.Error("server::GetMasterImages::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(imagesResp); err != nil {
		h.logger.Error("server::GetMasterImages::Write", err)
		return
	}

	h.logger.Info("Response sent")
}
