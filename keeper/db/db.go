package db

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/asxcandrew/wbrkev/keeper/config"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var db *pg.DB

// Customer db structure
type Customer struct {
	ID           uint64
	Name         string `validate:"required"`
	Email        string `validate:"required"`
	MobileNumber string `validate:"required"`
	CreatedAt    time.Time
}

func DB() *pg.DB {
	if db == nil {
		options := pg.Options{
			Addr:     fmt.Sprintf("%s:%s", config.C.DBHost, config.C.DBPort),
			User:     config.C.DBUser,
			Password: config.C.DBPass,
			Database: config.C.DBName,
		}
		if config.C.TLS {
			options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
		db = pg.Connect(&options)

		createSchema(db)
	}
	return db
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Customer)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
