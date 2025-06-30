// Defines Queries/Db operations related to Product Queries/Operations
package queries

import (
	"context"
	"lapbytes/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertLaptop(pool *pgxpool.Pool, lp model.Laptop) (product_id int, err error) {

	stmt := `
	INSERT INTO products (name, brand, operatingsystem, operatingsystemversion, 
    hdd, ssd, hddsize, ssdsize, ramsize, 
    cpumaker, cpugen, cpumodel, yom, imageurl, price, screensize,
    hasgpu, gpumake, gpumaker, hasigpu, isinstock, instock)
	
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
	
	RETURNING id
	`

	err = pool.QueryRow(context.Background(), stmt,
		lp.Name,
		lp.Brand,
		lp.Operating_system,
		lp.Operating_system_version,
		lp.HDD,
		lp.SSD,
		lp.HDD_size,
		lp.SSD_size,
		lp.Ram_size,
		lp.CPU_maker,
		lp.CPU_gen,
		lp.CPU_model,
		lp.YOM,
		lp.Image_url,
		lp.Price,
		lp.Screen_size,
		lp.Has_gpu,
		lp.Gpu_make,
		lp.Gpu_maker,
		lp.Has_igpu,
		lp.Is_in_stock,
		lp.In_stock,
	).Scan(&product_id)
	if err != nil {
		return 0, err
	}
	return product_id, nil

}

func QueryLaptop(pool *pgxpool.Pool, id int) (laptop model.Laptop, err error) {
	stmt := `
	SELECT id, name, brand, operatingsystem, operatingsystemversion, 
           hdd, ssd, hddsize, ssdsize, ramsize, 
           cpumaker, cpugen, cpumodel, yom, imageurl, price, screensize,
           hasgpu, gpumake, gpumaker, hasigpu, isinstock
	FROM products WHERE id=$1
`

	err = pool.QueryRow(context.Background(), stmt, id).Scan(

		&laptop.Id,
		&laptop.Name,
		&laptop.Brand,
		&laptop.Operating_system,
		&laptop.Operating_system_version,
		&laptop.HDD,
		&laptop.SSD,
		&laptop.HDD_size,
		&laptop.SSD_size,
		&laptop.Ram_size,
		&laptop.CPU_maker,
		&laptop.CPU_gen,
		&laptop.CPU_model,
		&laptop.YOM,
		&laptop.Image_url,
		&laptop.Price,
		&laptop.Screen_size,
		&laptop.Has_gpu,
		&laptop.Gpu_make,
		&laptop.Gpu_maker,
		&laptop.Has_igpu,
		&laptop.Is_in_stock,
	)
	if err != nil {
		return model.Laptop{}, err
	}
	return laptop, nil
}

func DeleteLaptop(pool *pgxpool.Pool, id int) error {

	stmt := `
		DELETE FROM products WHERE id=$1
	`
	_, err := pool.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	return nil
}
