package mysql

import (
	"backend/utils/config"
	"backend/utils/log"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbItem  *dbPgsql
	dbAdmin *dbPgsql
)

type dbPgsql struct {
	dbPq *sql.DB
}

func InitDBConnection() error {
	err := InitConnectionItem()
	if err != nil {
		return err
	}

	err = InitConnectionAdmin()
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

func InitConnectionAdmin() error {
	log.Log.Println("start init DB Admin", nil)
	dbAdmin = new(dbPgsql)
	conn, errdb := sql.Open("mysql", config.MyConfig.DbAdmin)
	if errdb != nil {
		return errdb
	}
	if err := conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxOpenConns(3)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(5 * time.Minute)

	dbAdmin.dbPq = conn

	return nil
}

func GetConnectionItem() (*sql.DB, error) {
	return dbItem.GetConnection()
}

func GetConnectionAdmin() (*sql.DB, error) {
	return dbAdmin.GetConnection()
}

func (dpq *dbPgsql) GetConnection() (*sql.DB, error) {
	if err := dpq.dbPq.Ping(); err != nil {
		return nil, err
	}
	return dpq.dbPq, nil
}

func CloseDBConnection() {
	closeItemConnection()
	closeAdminConnection()
}

func closeItemConnection() {
	log.Log.Println("Closing Item DB connection")
	dbItem.dbPq.Close()
}

func closeAdminConnection() {
	log.Log.Println("Closing Admin DB connection")
	dbAdmin.dbPq.Close()
}
