package config

import (
	"encoding/json"
	"os"
)

func InitConf[T any](v *T) error {
	path := os.Getenv("CONF")
	if path == "" {
		path = "./server.conf"
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, v)
	if err != nil {
		return err
	}

	return nil
}
