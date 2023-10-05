package models

import "time"

type City struct {
	ID        uint      `gorm:"id"`
	CreatedAt time.Time `gorm:"created_at"`
	Name      string    `gorm:"name"`
}

type ServiceCategory struct {
	ID        uint      `gorm:"id"`
	CreatedAt time.Time `gorm:"created_at"`
	Name      string    `gorm:"name"`
}

type Service struct {
	ID        uint      `gorm:"id"`
	CreatedAt time.Time `gorm:"created_at"`
	Name      string    `gorm:"name"`
	CatID     uint      `gorm:"cat_id"`
	CatName   string    `gorm:"cat_name"`
}

type MasterServRelation struct {
	ID          uint      `gorm:"id"`
	CreatedAt   time.Time `gorm:"created_at"`
	MasterID    uint      `gorm:"master_id"`
	Name        string    `gorm:"name"`
	Description string    `gorm:"description"`
	Contact     string    `gorm:"contact"`
	CityID      uint      `gorm:"city_id"`
	CityName    string    `gorm:"city_name"`
	ServCatID   uint      `gorm:"serv_cat_id"`
	ServCatName string    `gorm:"serv_cat_name"`
	ServID      uint      `gorm:"serv_id"`
	ServName    string    `gorm:"serv_name"`
}

type MasterRegForm struct {
	ID          uint      `gorm:"id"`
	CreatedAt   time.Time `gorm:"created_at"`
	MasterID    uint      `gorm:"master_id"`
	Name        string    `gorm:"name"`
	Description string    `gorm:"description"`
	Contact     string    `gorm:"contact"`
	CityID      uint      `gorm:"city_id"`
	CityName    string    `gorm:"city_name"`
	ServCatID   uint      `gorm:"serv_cat_id"`
	ServCatName string    `gorm:"serv_cat_name"`
	ServID      uint      `gorm:"serv_id"`
	ServName    string    `gorm:"serv_name"`
}

type MasterImages struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;notNull"`
	MasterID uint   `gorm:"master_id"`
	URL      string `gorm:"url"`
}
