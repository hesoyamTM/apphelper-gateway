package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string  `yaml:"env" env-required:"true" env:"ENV"`
	PublicKey   string  `yaml:"public_key" env-required:"true" env:"PUBLIC_KEY"`
	Http        HTTP    `yaml:"http"`
	GrpcClients Clients `yaml:"clients"`
}

type HTTP struct {
	Address     string        `yaml:"address" env-required:"true" env:"HTTP_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true" env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true" env:"HTTP_IDLE_TIMEOUT"`
}

type Clients struct {
	SsoAddress      string `yaml:"sso_address" env-required:"true" env:"SSO_ADDRESS"`
	ReportAddress   string `yaml:"report_address" env-required:"true" env:"REPORT_ADDRESS"`
	ScheduleAddress string `yaml:"schedule_address" env-required:"true" env:"SCHEDULE_ADDRESS"`
}

func fetchConfigPath() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "config path")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("CONFIG_PATH")
	}

	return cfgPath
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env")
	}

	cfgPath := fetchConfigPath()
	if cfgPath != "" {
		return MustLoadByPath(cfgPath)
	}

	return MustLoadEnv()
}

func MustLoadEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func MustLoadByPath(cfgPath string) *Config {
	if cfgPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic("config file does not exist: " + err.Error())
		}

		panic(err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
