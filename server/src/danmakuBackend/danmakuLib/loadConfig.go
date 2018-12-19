package danmakuLib

import (
	//"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type ConfigObject struct {
	Port int
	Database string
	DBserver string
	DBport int
	DBuser string
	DBpassword string
	DBsource string
	WSToken string
}
// external accessible function / var have to start with uppercase letter.

func GetConfig() ConfigObject {
	configFile, err := ioutil.ReadFile("./config/config.json")

	if err != nil{
		fmt.Print("Failed to read config file. Exit now.\n")
		os.Exit(1)
	}

	var config ConfigObject
	err = json.Unmarshal(configFile, &config)
	if err != nil{
		fmt.Print("failed to parse json config. Exit now.\n ")
		os.Exit(0)
	}

	return config
}
