package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // this is because sqlx package needs pq driver
)

// DBConf is a DB configuration struct
type DBConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func getConnString(conf DBConf) string {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBName)
	return connString
}

func mustGetDB(conf DBConf) *sqlx.DB {
	connString := getConnString(conf)
	db := sqlx.MustConnect("postgres", connString)
	return db
}

func loadConf() DBConf {
	return DBConf{Host: "localhost", Port: 5432, User: "postgres", Password: "postgres", DBName: "postgres"}
}

// MustGetDb gets you configured db connection
func MustGetDb() *sqlx.DB {
	conf := loadConf()
	return mustGetDB(conf)
}

// func insertURI(db *sqlx.DB, ) string {

// }
