// //Defines Queries/Db operations related to users Operations
package queries

import (
	"context"
	"lapbytes/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertUser(pool *pgxpool.Pool, user model.User) (userId int, err error) {
	//Lojik goes here
	stmt := `
	INSERT INTO users (username,email,passwordhash,isadmin,accesslevel,createdat,updatedat)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`
	err = pool.QueryRow(context.Background(), stmt,
		user.Username,
		user.Email,
		user.Password_hash,
		user.Is_admin,
		user.Access_level,
		user.Created_at,
		user.Updated_at).Scan(&userId)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func GetUserHash(pool *pgxpool.Pool, email string) (string, error) {
	//Sanitize before bringing it here
	var passwordhash string
	stmt := `
	SELECT passwordhash FROM users WHERE email=$1
	`
	err := pool.QueryRow(context.Background(), stmt, email).Scan(&passwordhash)
	if err != nil {
		return "", err
	}

	return passwordhash, nil

}
