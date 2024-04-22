package dbadapter

import (
	"bot/internal/config"
	"bot/internal/entities"
	"bot/internal/mapper"
	"bot/internal/models"
	"fmt"
	"time"

	"bot/internal/logger"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBAdapter struct {
	logger logger.Logger
	cfg    *config.Config
	DBConn *gorm.DB
}

func NewDbAdapter(logger logger.Logger, cfg *config.Config) (*DBAdapter, error) {

	psqlconf := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PsqlHost,
		cfg.PsqlPort,
		cfg.PsqlUser,
		cfg.PsqlPass,
		cfg.PsqlDb,
	)

	DBConn, err := gorm.Open(postgres.Open(psqlconf), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DBAdapter{logger: logger, cfg: cfg, DBConn: DBConn}, nil
}

func (d *DBAdapter) AutoMigrate() error {
	if err := d.DBConn.AutoMigrate(&models.City{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.ServiceCategory{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.Service{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.MasterServRelation{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.MasterRegForm{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.MasterImages{}); err != nil {
		return err
	}
	d.logger.Info("Auto-migration: success")
	return nil
}

func (d *DBAdapter) GetCities(servID string, page, limit int) ([]*entities.City, error) {

	if len(servID) != 0 {
		return d.GetCitiesByService(servID, page, limit)
	}

	cities := make([]*models.City, 0)
	if err := d.DBConn.Offset(page * limit).Limit(limit).Find(&cities).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.City, 0)
	for _, city := range cities {
		result = append(result, mapper.FromCityModel(city))
	}

	return result, nil
}

func (d *DBAdapter) GetCitiesByService(servID string, page, limit int) ([]*entities.City, error) {

	masterServRelations := make([]*models.MasterServRelation, 0)
	query := d.DBConn.Offset(page * limit).Limit(limit)
	query = query.Where("serv_id = ?", servID).Select("DISTINCT ON (city_id) city_id, city_name")
	if err := query.Find(&masterServRelations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.City, 0)
	for _, relation := range masterServRelations {
		result = append(result, &entities.City{
			ID:   relation.CityID,
			Name: relation.CityName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServCategories(cityID string, page, limit int) ([]*entities.ServiceCategory, error) {

	if len(cityID) != 0 {
		return d.GetServCategoriesByCity(cityID, page, limit)
	}

	categories := make([]*models.ServiceCategory, 0)
	if err := d.DBConn.Offset(page * limit).Limit(limit).Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.ServiceCategory, 0)
	for _, category := range categories {
		result = append(result, mapper.FromServCatModel(category))
	}

	return result, nil
}

func (d *DBAdapter) GetServCategoriesByCity(cityID string, page, limit int) ([]*entities.ServiceCategory, error) {

	masterServRelations := make([]*models.MasterServRelation, 0)
	query := d.DBConn.Offset(page * limit).Limit(limit)
	query = query.Where("city_id = ?", cityID).Select("DISTINCT ON (serv_cat_id) serv_cat_id, serv_cat_name")
	if err := query.Find(&masterServRelations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.ServiceCategory, 0)
	for _, relation := range masterServRelations {
		result = append(result, &entities.ServiceCategory{
			ID:   relation.ServCatID,
			Name: relation.ServCatName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServices(categoryID, cityID string, page, limit int) ([]*entities.Service, error) {

	if len(cityID) != 0 {
		return d.GetServicesByCity(categoryID, cityID, page, limit)
	}

	return d.GetServicesByCategory(categoryID, page, limit)
}

func (d *DBAdapter) GetServicesByCity(categoryID, cityID string, page, limit int) ([]*entities.Service, error) {

	masterServRelations := make([]*models.MasterServRelation, 0)
	query := d.DBConn.Offset(page * limit).Limit(limit)
	if len(categoryID) != 0 {
		query = query.Where("serv_cat_id = ?", categoryID)
	}

	query = query.Where("city_id = ?", cityID).Select("DISTINCT ON (serv_id) serv_id, serv_name, serv_cat_id, serv_cat_name")
	if err := query.Find(&masterServRelations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Service, 0)
	for _, relation := range masterServRelations {
		result = append(result, &entities.Service{
			ID:      relation.ServID,
			Name:    relation.ServName,
			CatID:   relation.ServCatID,
			CatName: relation.ServCatName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServicesByCategory(categoryID string, page, limit int) ([]*entities.Service, error) {

	query := d.DBConn.Offset(page * limit).Limit(limit)
	if len(categoryID) != 0 {
		query = query.Where("cat_id = ?", categoryID)
	}

	services := make([]*models.Service, 0)
	if err := query.Find(&services).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Service, 0)
	for _, service := range services {
		result = append(result, mapper.FromServiceModel(service))
	}

	return result, nil
}

func (d *DBAdapter) GetMasters(cityID, servCatID, servID string, page, limit int) ([]*entities.Master, error) {

	query := d.DBConn.Offset(page * limit).Limit(limit)
	if len(cityID) != 0 {
		query = query.Where("city_id = ?", cityID)
	}
	if len(servCatID) != 0 {
		query = query.Where("serv_cat_id = ?", servCatID)
	}
	if len(servID) != 0 {
		query = query.Where("serv_id = ?", servID)
	}

	masters := make([]*models.MasterServRelation, 0)
	if err := query.Select("DISTINCT ON (master_id) *").Find(&masters).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Master, 0)
	for _, master := range masters {
		images, err := d.GetMasterImages(master.MasterID)
		if err != nil {
			d.logger.Errorf("Failed to find image for %d %s %s", master.MasterID, master.Name, err)
		}
		result = append(result, mapper.FromMasterServRelationModel(master, images))
	}
	return result, nil
}

func (d *DBAdapter) GetMaster(masterID string) (*entities.Master, error) {

	return nil, nil
}

func (d *DBAdapter) GetMasterImages(masterID string) ([]string, error) {

	imgRecs := make([]*models.MasterImages, 0)
	if err := d.DBConn.Where("master_id = ?", masterID).Find(&imgRecs).Error; err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, rec := range imgRecs {
		result = append(result, fmt.Sprintf("%s/%s/%s", d.cfg.ImagePrefix, masterID, rec.Name))
	}

	return result, nil
}

func (d *DBAdapter) SaveCity(name string) (string, error) {
	id := uuid.NewString()
	city := &models.City{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      name,
	}
	if err := d.DBConn.Create(city).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New city added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveServiceCategory(name string) (string, error) {
	id := uuid.NewString()
	service := &models.ServiceCategory{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      name,
	}
	if err := d.DBConn.Create(service).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New service category added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveService(name, categoryID string) (string, error) {
	id := uuid.NewString()

	category := models.ServiceCategory{}
	if err := d.DBConn.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return "", err
	}

	service := &models.Service{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      name,
		CatID:     category.ID,
		CatName:   category.Name,
	}
	if err := d.DBConn.Create(service).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New service added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveMasterRegForm(master *entities.MasterRegForm) (string, error) {
	city := models.City{}
	if err := d.DBConn.Where("id = ?", master.CityID).First(&city).Error; err != nil {
		return "", err
	}

	id := uuid.NewString()
	regForm := models.MasterRegForm{
		CreatedAt:   time.Now(),
		MasterID:    id,
		Name:        master.Name,
		Contact:     master.Contact,
		Description: master.Description,
		CityID:      city.ID,
		CityName:    city.Name,
	}

	forms := make([]models.MasterRegForm, 0)
	for _, servID := range master.ServIDs {
		service := models.Service{}
		if err := d.DBConn.Where("id = ?", servID).First(&service).Error; err != nil {
			return "", err
		}
		regForm.ID = uuid.NewString()
		regForm.ServCatID = service.CatID
		regForm.ServCatName = service.CatName
		regForm.ServID = service.ID
		regForm.ServName = service.Name
		forms = append(forms, regForm)
	}

	if err := d.DBConn.Create(&forms).Error; err != nil {
		return "", err
	}
	d.logger.Infof("Form saved successfully, id: %s, name: %s", id, master.Name)
	return id, nil
}

func (d *DBAdapter) SaveMaster(id string) (string, error) {

	masters := make([]*models.MasterRegForm, 0)
	if err := d.DBConn.Where("master_id = ?", id).Find(&masters).Error; err != nil {
		return "", err
	}

	result := make([]*models.MasterServRelation, 0)
	for _, master := range masters {
		result = append(result, &models.MasterServRelation{
			ID:          master.ID,
			CreatedAt:   time.Now(),
			MasterID:    master.MasterID,
			Name:        master.Name,
			Description: master.Description,
			Contact:     master.Contact,
			CityID:      master.CityID,
			CityName:    master.CityName,
			ServCatID:   master.ServCatID,
			ServCatName: master.ServCatName,
			ServID:      master.ServID,
			ServName:    master.ServName,
		})
	}

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Create(&result).Error; err != nil {
		return "", err
	}

	if err := tx.Where("master_id = ?", id).Delete(&models.MasterRegForm{}).Error; err != nil {
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	d.logger.Infof("New master added successfully, id: %s", id)
	return id, nil
}

func (d *DBAdapter) SaveMasterImage(masterID, name string) error {

	imgRec := models.MasterImages{
		MasterID: masterID,
		Name:     name,
	}

	if err := d.DBConn.Create(&imgRec).Error; err != nil {
		return err
	}

	d.logger.Infof("Image saved successfully, %s", name)
	return nil
}

func (d *DBAdapter) UpdateCity(city *entities.City) error {

	update := models.City{
		ID:        city.ID,
		CreatedAt: time.Now(),
		Name:      city.Name,
	}

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Save(&update).Error; err != nil {
		return err
	}

	query := tx.Model(&models.MasterServRelation{}).Where("city_id = ?", city.ID)
	if err := query.UpdateColumn("city_name", city.Name).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("City changed successfully: %s", city.Name)
	return nil
}

func (d *DBAdapter) UpdateServCategory(category *entities.ServiceCategory) error {

	update := models.ServiceCategory{
		ID:        category.ID,
		CreatedAt: time.Now(),
		Name:      category.Name,
	}

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Save(&update).Error; err != nil {
		return err
	}

	query := tx.Model(&models.Service{}).Where("cat_id = ?", category.ID)
	if err := query.UpdateColumn("cat_name", category.Name).Error; err != nil {
		return err
	}

	query = tx.Model(&models.MasterServRelation{}).Where("serv_cat_id = ?", category.ID)
	if err := query.UpdateColumn("serv_cat_name", category.Name).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("Service category changed successfully: %s", category.Name)
	return nil
}

func (d *DBAdapter) UpdateService(service *entities.Service) error {

	category := models.ServiceCategory{}
	if err := d.DBConn.Where("id = ?", service.CatID).First(&category).Error; err != nil {
		return err
	}

	update := models.Service{
		ID:        service.ID,
		CreatedAt: time.Now(),
		CatID:     category.ID,
		CatName:   category.Name,
		Name:      service.Name,
	}

	relation := models.MasterServRelation{
		ServCatID:   update.CatID,
		ServCatName: update.CatName,
		ServName:    update.Name,
	}

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Save(&update).Error; err != nil {
		return err
	}

	query := tx.Model(&models.MasterServRelation{}).Where("serv_id = ?", service.ID)
	if err := query.UpdateColumns(&relation).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("Service changed successfully: %s", service.Name)
	return nil
}

func (d *DBAdapter) UpdateMaster() {

}

func (d *DBAdapter) DeleteCity(id string) error {
	if err := d.DBConn.Where("id = ?", id).Delete(&models.City{}).Error; err != nil {
		return err
	}

	d.logger.Infof("City was deleted successfully: %s", id)
	return nil
}

func (d *DBAdapter) DeleteServCategory(id string) error {

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Where("id = ?", id).Delete(&models.ServiceCategory{}).Error; err != nil {
		return err
	}

	if err := tx.Where("cat_id = ?", id).Delete(&models.Service{}).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("ServiceCategory was deleted successfully: %s", id)
	return nil
}

func (d *DBAdapter) DeleteService(id string) error {
	if err := d.DBConn.Where("id = ?", id).Delete(&models.Service{}).Error; err != nil {
		return err
	}

	d.logger.Infof("Service was deleted successfully: %s", id)
	return nil
}

func (d *DBAdapter) DeleteMaster(id string) error {

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Where("master_id = ?", id).Delete(&models.MasterServRelation{}).Error; err != nil {
		return err
	}

	if err := tx.Where("master_id = ?", id).Delete(&models.MasterImages{}).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("Master was deleted successfully: %s", id)
	return nil
}
