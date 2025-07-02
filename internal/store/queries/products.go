// Defines Queries/Db operations related to Product Queries/Operations
package queries

import (
	"context"
	"lapbytes/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func QueryLaptops(pool *pgxpool.Pool, limit int, offset int) (laptops []model.Laptop, err error) {
	stmt := `
		SELECT id, name, brand, operatingsystem, operatingsystemversion, 
           hdd, ssd, hddsize, ssdsize, ramsize, 
           cpumaker, cpugen, cpumodel, yom, imageurl, price, screensize,
           hasgpu, gpumake, gpumaker, hasigpu, isinstock
		FROM products 
		ORDER BY createdat DESC
		LIMIT $1 OFFSET $2
`
	rows, err := pool.Query(context.Background(), stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p model.Laptop
		//if err occurs on scan, return incomplete data*
		err = rows.Scan(&p.Id,
			&p.Name,
			&p.Brand,
			&p.Operating_system,
			&p.Operating_system_version,
			&p.HDD,
			&p.SSD,
			&p.HDD_size,
			&p.SSD_size,
			&p.Ram_size,
			&p.CPU_maker,
			&p.CPU_gen,
			&p.CPU_model,
			&p.YOM,
			&p.Image_url,
			&p.Price,
			&p.Screen_size,
			&p.Has_gpu,
			&p.Gpu_make,
			&p.Gpu_maker,
			&p.Has_igpu,
			&p.Is_in_stock,
		)
		if err != nil {
			return nil, err
		}
		laptops = append(laptops, p)
	}
	return laptops, nil
}
