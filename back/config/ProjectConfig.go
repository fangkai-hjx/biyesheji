package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type projectConfig struct {
	ProjectName string `json:"project_name"`
	ProjectMode string `json:"project_mode"`
	ProjectHost string `json:"project_host"`
	ProjectPort string `json:"project_port"`

	K8sEvn string `json:"k8s_evn"`

	HarborUrl      string `json:"harbor_url"`
	HarborUsername string `json:"harbor_username"`
	HarborPassword string `json:"harbor_password"`

	DriverName   string `json:"driver_name"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Charset      string `json:"charset"`

	RedisUrl     string `json:"redis_url"`
	RedisDB      int `json:"redis_db"`

	PrometheusUrl string `json:"prometheus_url"`

}

var (
	ProjectConfig *projectConfig
)

func init() {
	ParseConfig()
}
func ParseConfig() (*projectConfig, error) {

	file, err := os.Open("./config/project.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&ProjectConfig); err != nil {
		return nil, err
	}
	return ProjectConfig, err
}
func (config *projectConfig) GetDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		config.UserName, config.Password, config.Host, config.Port, config.DatabaseName, config.Charset)
}