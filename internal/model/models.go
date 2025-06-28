package model

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Laptop struct {
	//Properties
	Id                       uuid.UUID
	Name                     string
	Brand                    string
	Operating_system         string
	Operating_system_version string
	HDD                      bool
	SSD                      bool
	HDD_size                 float64
	SSD_size                 float64
	Ram_size                 float64
	CPU_maker                string
	CPU_gen                  string
	CPU_model                string
	YOM                      string
	Image_url                string
	Price                    float64
	Screen_size              float64
	Has_gpu                  bool
	Gpu_make                 string
	Gpu_maker                string
	Has_igpu                 string

	//Stock
	Is_in_stock bool
	In_stock    int
	Created_at  time.Time
	Updated_at  time.Time
}

type User struct {
	Id            uuid.UUID
	Username      string
	Email         string
	Password_hash string
	Is_admin      bool
	Access_level  string
	Created_at    time.Time
	Updated_at    time.Time
}

type Cart struct {
}

type LoginResponse struct {
	AccessToken string `json:access_token`
	TokenType   string `json:token_type`
}

type RefreshHttpOnlyCookie struct {
	Name     string
	Value    string
	Path     string
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
	Expires  time.Time
}
