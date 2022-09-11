package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BurntSushi/toml"
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
	// 读取环境变量中的配置文件类型
	confType := os.Getenv("CONF_TYPE")
	if confType == "" {
		confType = "json"
	}

	switch confType {
	case "json":
		err = json.Unmarshal(file, v)
	case "toml":
		err = toml.Unmarshal(file, v)
	default:
		// 不支持的配置文件类型
		log.Fatalf("unsupported config file type<%s>", confType)
	}
	if err != nil {
		return err
	}

	return nil
}
