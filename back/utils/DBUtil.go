package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"t/back/config"
)

var (
	DB *sqlx.DB = nil
)

func GetDBClient()*sqlx.DB {
	if DB != nil{
		return DB
	}
	db, err := sqlx.Open(config.DBConfig.DriverName, config.DBConfig.GetDataSourceName())
	if err != nil {
		log.Printf("DBUtil.getDBClient.Open: %s\n", err)
		return nil
	}
	DB = db
	return DB
}
