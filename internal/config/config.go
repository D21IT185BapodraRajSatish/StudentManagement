package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func Mustload() *Config {

	configpath := os.Getenv("CONFIG_PATH")

	if configpath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configpath := *flags

		if configpath == "" {
			log.Fatal("Config Path is not set")
		}
	}

	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist %s", configpath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configpath, &cfg)

	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}
	return &cfg
}
