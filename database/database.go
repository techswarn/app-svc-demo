package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
	"fmt"
)

type Config struct {
	user string
	password string
	proto string
	hostname string
	dbname string
}

type Country struct {
	Code string
	Name string
	Continent string
	Region string
	Area float64
	GovernmentFrom string
}

type DB struct {
	db *sql.DB
}

func NewCon() (*DB, error) {

	config := &Config{
		user: os.Getenv("DATABASE_USER"),
	    password: os.Getenv("DATABASE_PASSWORD"),
	    proto: "tcp",
	    hostname: os.Getenv("DATABASE_ADDR"),
	    dbname: os.Getenv("DATABASE_DBNAME"),
	}

	fmt.Printf("db config: %#v", config)

	cfg := mysql.NewConfig()
    cfg.User = config.user
    cfg.Passwd = config.password
    cfg.Net = config.proto
    cfg.Addr = config.hostname
    cfg.DBName = config.dbname

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        fmt.Printf("Error connecting to mysql database : %s", err)
		return nil, err
    }

    pingErr := db.Ping()
    if pingErr != nil {
        fmt.Printf("Database ping failed %s", pingErr)
		return nil, pingErr
    }
    fmt.Println("Connected!")

	return &DB{
		db: db,
	}, nil
}

func (d *DB) close() {
	err := d.db.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DB) GetCountries() ([]Country, error) {
	var countries []Country
	rows, err := d.db.Query("SELECT Code, Name, Continent, Region, SurfaceArea, GovernmentForm FROM country ORDER BY Name")
	if err != nil {
		fmt.Printf("Error while query data %s", err)
		return nil, err
	}
	defer rows.Close()
    for rows.Next() {
		var cnty Country
		err := rows.Scan(&cnty.Code, &cnty.Name, &cnty.Continent, &cnty.Region, &cnty.Area, &cnty.GovernmentFrom)
		if err != nil {
			fmt.Printf("Error while scanning %s", err)
			return nil, err
		}
		countries = append(countries, cnty)
	}
	return countries, nil
}
