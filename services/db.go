package services

import "database/sql"

var defaultDB *sql.DB

func SetDB(db *sql.DB) {
	defaultDB = db
}

func GetDB() *sql.DB {
	return defaultDB
}
