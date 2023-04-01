package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := "host=localhost " +
		"port=5432 " +
		"user=postgres " +
		"password=testperson " +
		"dbname=person " +
		"TimeZone=Europe/Moscow"
		
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
