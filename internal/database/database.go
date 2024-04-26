package database

import (
	"effective_mobile_2/internal/config"
	"effective_mobile_2/internal/repository/gorm/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
}

var db Database

func Connect() error {
	var err error

	db.Gorm, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Cfg().Postgres.Url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return err
	}

	return nil
}

func Migrate() error {
	err := db.Gorm.AutoMigrate(
		&entity.People{},
		&entity.Car{},
	)

	if err != nil {
		return err
	}

	return nil
}

func Db() *Database {
	return &db
}
