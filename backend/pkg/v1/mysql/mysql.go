package mysql

import (
	"backend/utils/config"
	"backend/utils/log"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbItem *dbPgsql
	dbAuth *dbPgsql
)

type dbPgsql struct {
	dbPq *sql.DB
}

func InitDBConnection() error {
	err := InitConnectionItem()
	if err != nil {
		return err
	}

	err = InitConnectionAuth()
	if err != nil {
		return err
	}

	return nil
}

func InitConnectionItem() error {
	log.Log.Println("start init DB Item", nil)
	dbItem = new(dbPgsql)
	conn, errdb := sql.Open("mysql", config.MyConfig.DbItem)
	if errdb != nil {
		return errdb
	}
	if err := conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxOpenConns(3)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(5 * time.Minute)

	dbItem.dbPq = conn

	return nil
}

func InitConnectionAuth() error {
	log.Log.Println("start init DB User", nil)
	dbAuth = new(dbPgsql)
	conn, errdb := sql.Open("mysql", config.MyConfig.DbAuth)
	if errdb != nil {
		return errdb
	}
	if err := conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxOpenConns(3)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(5 * time.Minute)

	dbAuth.dbPq = conn

	return nil
}

func GetConnectionItem() (*sql.DB, error) {
	return dbItem.GetConnection()
}

func GetConnectionUser() (*sql.DB, error) {
	return dbAuth.GetConnection()
}

func (dpq *dbPgsql) GetConnection() (*sql.DB, error) {
	if err := dpq.dbPq.Ping(); err != nil {
		return nil, err
	}
	return dpq.dbPq, nil
}

func CloseDBConnection() {
	closeItemConnection()
	closeUserConnection()
}

func closeItemConnection() {
	log.Log.Println("Closing Item DB connection")
	dbItem.dbPq.Close()
}

func closeUserConnection() {
	log.Log.Println("Closing User DB connection")
	dbAuth.dbPq.Close()
}
