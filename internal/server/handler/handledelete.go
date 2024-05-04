package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
		h.logger.Errorf("server::DeleteCity::DeleteCity: %s", err.Error())
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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)

	if err := h.DBAdapter.DeleteServCategory(params["category_id"]); err != nil {
		h.logger.Errorf("server::DeleteServCategory::DeleteServCategory: %s", err.Error())
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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)

	if err := h.DBAdapter.DeleteService(params["service_id"]); err != nil {
		h.logger.Errorf("server::DeleteService::DeleteService: %s", err.Error())
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
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	if err := h.DBAdapter.DeleteMaster(masterID); err != nil {
		h.logger.Errorf("server::DeleteMaster::DeleteMaster: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.MinIOAdapter.DeleteMasterImages(masterID); err != nil {
		h.logger.Errorf("server::DeleteMaster::DeleteBucket: %s", err.Error())
	}

	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}

// @Summary Delete master image
// @Description Delete an image of a master from the system
// @Tags Master
// @Param master_id path string true "ID of the master"
// @Param image_name path string true "Name of the image"
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {string} string "Error message"
// @Router /masters/{master_id}/images/{image_name} [delete]
func (h *Handler) DeleteMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s %s", req.Method, req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]
	imageName := params["image_name"]

	if err := h.MinIOAdapter.DeleteMasterImage(masterID, imageName); err != nil {
		h.logger.Errorf("server::DeleteMasterImage::DeleteMasterImage: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}
