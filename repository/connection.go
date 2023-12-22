package repository

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func DbConnection() (*sqlx.DB, error) {
	connection := fmt.Sprintf("%s:%s@(%s:%s)/%s", "nitesh", "password", "localhost", "3306", "splitwise")

	db, err := sqlx.Connect("mysql", connection)
	if err != nil {
		log.Println("error in establishing DB connection", err.Error())
		return nil, err
	}

	return db.Unsafe(), nil

}
