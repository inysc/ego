package config

import (
	"encoding/json"
	"os"
)

var srvName string

func InitConf[T any](v *T) error {
	file, err := os.ReadFile(os.Getenv("CONF"))
	if err != nil {
		return err
	}

	return json.Unmarshal(file, v)
}

func SetSrvName(name string) { srvName = name }

func SrvName() string { return srvName }
