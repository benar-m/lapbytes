package model

import (
	"database/sql"
	"net/http"
	"time"
)

type Laptop struct {
	//Properties
	Id                       int            `json:"id" db:"id"`
	Name                     string         `json:"name" db:"name"`
	Brand                    string         `json:"brand" db:"brand"`
	Operating_system         string         `json:"operating_system" db:"operatingsystem"`
	Operating_system_version string         `json:"operating_system_version" db:"operatingsystemversion"`
	HDD                      bool           `json:"hdd" db:"hdd"`
	SSD                      bool           `json:"ssd" db:"ssd"`
	HDD_size                 float64        `json:"hdd_size" db:"hddsize"`
	SSD_size                 float64        `json:"ssd_size" db:"ssdsize"`
	Ram_size                 float64        `json:"ram_size" db:"ramsize"`
	CPU_maker                string         `json:"cpu_maker" db:"cpumaker"`
	CPU_gen                  string         `json:"cpu_generation" db:"cpugen"`
	CPU_model                string         `json:"cpu_model" db:"cpumodel"`
	YOM                      string         `json:"year_of_manufacture" db:"yom"`
	Image_url                string         `json:"image_url" db:"imageurl"`
	Price                    float64        `json:"price" db:"price"`
	Screen_size              float64        `json:"screen_size" db:"screensize"`
	Has_gpu                  bool           `json:"has_gpu" db:"hasgpu"`
	Gpu_make                 sql.NullString `json:"gpu_model" db:"gpumake"`
	Gpu_maker                sql.NullString `json:"gpu_manufacturer" db:"gpumaker"`
	Has_igpu                 bool           `json:"has_integrated_gpu" db:"hasigpu"`

	//Stock - Hidden from users
	Is_in_stock bool      `json:"is_in_stock" db:"isinstock"`
	In_stock    int       `json:"-" db:"instock"`
	Created_at  time.Time `json:"-" db:"createdat"`
	Updated_at  time.Time `json:"-" db:"updatedat"`
}
type User struct {
	Id            int       `db:"id"`
	Username      string    `db:"username"`
	Email         string    `db:"email"`
	Password_hash string    `db:"passwordhash"`
	Is_admin      bool      `db:"isadmin"`
	Access_level  int       `db:"accesslevel"` //There will be 5 levels of access with 0 being the highest (super user) and 4 the lowest
	Created_at    time.Time `db:"createdat"`
	Updated_at    time.Time `db:"updatedat"`
}

type Cart struct {
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
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
