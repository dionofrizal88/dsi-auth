package db

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

// NewPostgresDBConnection will create a PostgresSQL DB connection.
func NewPostgresDBConnection(conf config.Configuration) (*gorm.DB, error) {
	connectionString := "host=" + conf.DBHost +
		" port=" + conf.DBPort +
		" user=" + conf.DBUsername +
		" password=" + conf.DBPassword +
		" dbname=" + conf.DBName +
		" sslmode=disable"

	gormDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// NewPostgresDBTestConnection will create a PostgresSQL DB test connection.
func NewPostgresDBTestConnection(conf config.Configuration) (*gorm.DB, error) {
	connectionString := "host=" + conf.DBHost +
		" port=" + conf.DBPort +
		" user=" + conf.DBUsername +
		" password=" + conf.DBPassword +
		" dbname=" + conf.TestDBName +
		" sslmode=disable"

	gormDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
