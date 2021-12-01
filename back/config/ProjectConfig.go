package config

import (
	"bufio"
	"encoding/json"
	"os"
)


type projectConfig struct {
	ProjectName     string `json:"project_name"`
	ProjectMode     string `json:"project_mode"`
	ProjectHost     string `json:"project_host"`
	ProjectPort     string `json:"project_port"`
}

var (
	ProjectConfig *projectConfig
)

func init()  {
	ParseConfig()
}
func ParseConfig() (*projectConfig, error) {

	file, err := os.Open("./back/config/project.json")
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
