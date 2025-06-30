CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    operatingsystem VARCHAR(255) NOT NULL,
    operatingsystemversion VARCHAR(255) NOT NULL,
    hdd BOOLEAN NOT NULL DEFAULT FALSE,
    ssd BOOLEAN NOT NULL DEFAULT FALSE,
    hddsize FLOAT NOT NULL DEFAULT 0.0,
    ssdsize FLOAT NOT NULL DEFAULT 0.0,
    ramsize FLOAT NOT NULL,
    cpumaker VARCHAR(255) NOT NULL,
    cpugen VARCHAR(255) NOT NULL,
    cpumodel VARCHAR(255) NOT NULL,
    yom VARCHAR(255) NOT NULL,
    imageurl VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    screensize FLOAT NOT NULL,
    hasgpu BOOLEAN NOT NULL DEFAULT FALSE,
    gpumake VARCHAR(255) ,
    gpumaker VARCHAR(255) ,
    hasigpu BOOLEAN NOT NULL DEFAULT FALSE,
    isinstock BOOLEAN NOT NULL DEFAULT FALSE,
    instock INT NOT NULL DEFAULT 0,
    createdat TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedat TIMESTAMP NOT NULL DEFAULT NOW()

);
CREATE UNIQUE INDEX idx_products_name_price_os ON products (name, price, operatingsystem);



/*
type Laptop struct {
	//Properties
	Id                       int
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
	Has_igpu                 bool

	//Stock
	Is_in_stock bool
	In_stock    int
	Created_at  time.Time
	Updated_at  time.Time
}*/