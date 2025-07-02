// Defines Queries/Db operations related to users with admin access (Access Level 0)
package queries

import (
	"context"
	"errors"
	"fmt"
	"lapbytes/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
## Admin (Protected)
- `POST /api/admin/products` — Add new laptop
- `DELETE /api/admin/products/{id}` — Delete laptop
- `GET /api/admin/orders` — View all orders
- `GET /api/admin/users` — View all registered users
- `GET /api/admin/users/{id}` — View specific user details
*/

func InsertLaptop(pool *pgxpool.Pool, lp model.Laptop) (product_id int, err error) {

	stmt := `
	INSERT INTO products (name, brand, operatingsystem, operatingsystemversion, 
    hdd, ssd, hddsize, ssdsize, ramsize, 
    cpumaker, cpugen, cpumodel, yom, imageurl, price, screensize,
    hasgpu, gpumake, gpumaker, hasigpu, isinstock)
	
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
	
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
	).Scan(&product_id)
	if err != nil {
		return 0, err
	}
	return product_id, nil

}

func DeleteLaptop(pool *pgxpool.Pool, id int) error {

	stmt := `
		DELETE FROM products WHERE id=$1
	`
	result, err := pool.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}

func GetAllUsers(pool *pgxpool.Pool, limit, offset int) (users []model.User, err error) {
	stmt := `
	SELECT username,email,createdat,accesslevel
	FROM users
	ORDER BY createdat DESC
	LIMIT $1 OFFSET $2
	
	`
	rows, err := pool.Query(context.Background(), stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		if err = rows.Scan(
			&u.Username,
			&u.Email,
			&u.Created_at,
			&u.Access_level,
		); err != nil {
			return nil, err
		}
		users = append(users, u)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

}

func DeleteUser(pool *pgxpool.Pool, id int) error {
	stmt := `
		DELETE
		FROM users
		WHERE id=$1
	`
	result, err := pool.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user With id %d not found", id)
	}
	return nil

}

func GetUser(pool *pgxpool.Pool, id int) (user model.User, err error) {

	stmt := `
	SELECT username,email,createdat,accesslevel
	FROM users
	WHERE id=$1
	
	`
	err = pool.QueryRow(context.Background(), stmt, id).Scan(
		&user.Username,
		&user.Email,
		&user.Created_at,
		&user.Access_level)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return model.User{}, err
	}

	return user, nil //watch out
}
