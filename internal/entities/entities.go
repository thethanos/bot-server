package entities

const (
	PENDING = iota + 1
	APPROVED
	DECLINED
)

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ServiceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID      string `json:"id"`
	Name    string `json:"name" validate:"required"`
	CatID   string `json:"catID" validate:"required"`
	CatName string `json:"catName"`
}

type Image struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}

type Master struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description,omitempty"`
	Contact     string   `json:"contact" validate:"required"`
	CityID      string   `json:"cityID" validate:"required"`
	ServCatID   string   `json:"servCatID" validate:"required"`
	ServIDs     []string `json:"servIDs" validate:"required"`
	Status      uint     `json:"status" validate:"required"`
}

type MasterLong struct {
	Master
	ID string `json:"id" validate:"required"`
}

type MasterShort struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
	CityName    string   `json:"cityName"`
	ServCatName string   `json:"servCatName"`
	RegDate     string   `json:"regDate"`
	Images      []string `json:"images"`
}
