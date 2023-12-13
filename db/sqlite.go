package crawal_sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	// Kết nối đến cơ sở dữ liệu SQLite
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra xem bảng có tồn tại hay không, nếu không thì tạo mới
	createTable := `
		CREATE TABLE azur_lane (
			id_wallpaper INT PRIMARY KEY,
			file_name VARCHAR(255) NOT NULL,
			url VARCHAR(255) NOT NULL
		);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDb() *sql.DB {
	return db
}
