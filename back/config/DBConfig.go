package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type dbConfig struct {
	DriverName   string `json:"driver_name"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Charset      string `json:"charset"`
}

var (
	DBConfig   *dbConfig
	configPath = "./config/db.json"
)

func init() {
	ParseDbConfig()
	fmt.Println("mysql init successful...")
}

func ParseDbConfig() {
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("config.ParseDbConfig.Open: %s\n", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&DBConfig); err != nil {
		log.Printf("config.ParseDbConfig.Decode: %s\n", err)
		return
	}
}

func (config *dbConfig) GetDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		config.UserName, config.Password, config.Host, config.Port, config.DatabaseName, config.Charset)
}
