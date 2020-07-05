package dal

import (
	"database/sql"
	"fmt"
	"time"

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

type ShortenedURL struct {
	Id          int64          `db:"id"`
	LongUrl     string         `db:"long_url"`
	ShortUrl    string         `db:"short_url"`
	DateCreated time.Time      `db:"date_created"`
	HistCount   int            `db:"hit_count"`
	LastHit     sql.NullTime   `db:"last_hit"`
	UserAgent   sql.NullString `db:"user_agent"`
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
	return DBConf{Host: "localhost", Port: 5432, User: "postgres", Password: "postgres", DBName: "urls"}
}

// MustGetDb gets you configured db connection
func MustGetDb() *sqlx.DB {
	conf := loadConf()
	return mustGetDB(conf)
}

func GetOriginalURL(db *sqlx.DB, shortURL string) (string, bool) {
	longURLResult := ShortenedURL{}
	err := db.Get(&longURLResult, "SELECT * FROM urls WHERE short_url=$1", shortURL)
	if err != nil {
		return "", false
	}

	_, err = db.Exec("UPDATE urls SET hit_count = hit_count + 1, last_hit = now() WHERE short_url=$1", shortURL)
	if err != nil {
		return "", false
	}

	return longURLResult.LongUrl, true
}

func InsertShortenedURL(db *sqlx.DB, longURL string, userAgent string) string {
	id := 0
	err := db.QueryRow("INSERT INTO urls (long_url, user_agent)"+
		"VALUES ($1, $2) RETURNING id", longURL, userAgent).Scan(&id)
	if err != nil {
		panic(fmt.Sprintf("Cannot insert: %v", err))
	}
	shortened := ShortenedURL{}
	err = db.Get(&shortened, "SELECT * FROM urls WHERE id=$1", id)
	if err != nil {
		panic(fmt.Sprintf("Could not obtain newly created entity: %v", err))
	}
	return shortened.ShortUrl
}
