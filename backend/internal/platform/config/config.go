package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type (
	Env         string
	FrontendUrl string
)

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvStage Env = "stage"
	EnvProd  Env = "prod"
)

type Config struct {
	Env         Env
	Auth        Auth
	FrontendUrl FrontendUrl
	Database    Database
	Migrations  Migrations
	Cache       Cache
	HttpServer  HttpServer
	Mq          Mq
}

type Auth struct {
	Jwt struct {
		PrivateKeyPath string
		PublicKeyPath  string
	}
}

type Database struct {
	Dsn            string
	MigrationsPath *string
}

type Cache struct {
	Dsn string
}

type Mq struct {
	Dsn string
}

type HttpServer struct {
	Addr string
}

type Migrations struct {
	DBDSN string
	Path  string
}

func Load() Config {
	appenv := Env(os.Getenv("APP_ENV"))
	switch appenv {
	case EnvLocal, EnvDev, EnvStage, EnvProd:
	default:
		appenv = EnvLocal
	}

	envFilePath := flag.String("env", ".env", "path to env file")
	configFilePath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	if err := godotenv.Load(*envFilePath); err != nil {
		log.Println(".env file not loaded")
	}

	k := koanf.New(".")

	if err := k.Load(file.Provider(*configFilePath), yaml.Parser()); err != nil {
		log.Fatal(err.Error())
	}

	k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, "APP_")), "__", ".")

			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil)

	cfg := Config{Env: appenv}

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	return cfg
}
