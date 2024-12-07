package sqlite

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "../../repo.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Добавить CREATE IF NOT EXIST

	log.Println("Успешное подключение к базе данных")
	return db, nil
}
