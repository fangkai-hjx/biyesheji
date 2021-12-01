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

func getDBClient() {
	db, err := sqlx.Open(config.DBConfig.DriverName, config.DBConfig.GetDataSourceName())
	if err != nil {
		log.Printf("DBUtil.getDBClient.Open: %s\n", err)
		return
	}

	DB = db
}

func init() {
	getDBClient()
}
