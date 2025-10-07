package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	// connectionString := "root:root@tcp(127.0.0.1:3306)/" + dbName
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, dbPort, dbName)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		// panic(err)
		return nil, err
	}
	return db, nil
}