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

var appenv string

type Config struct {
	Env         string
	Auth        AuthConfig
	FrontendUrl string
	Database    struct {
		Default DatabaseConfig
	}
	Cache       CacheConfig
	HttpServer  HttpServerConfig
	Mq          struct {
		Default MqConfig
	}
}

type AuthConfig struct {
	Jwt struct {
		PrivateKeyPath string
		PublicKeyPath  string
	}
}

type DatabaseConfig struct {
	Dsn string
}

type CacheConfig struct {
	Dsn string
}

type MqConfig struct {
	Dsn string
}

type HttpServerConfig struct {
	Addr string
}

func Load() Config {
	appenv = os.Getenv("APP_ENV")
	if "" == appenv {
		appenv = "local"
	}

	envPath := flag.String("env", ".env", "path to env file")
	configPath := flag.String("config", "configs/config.yaml", "path to config")
	flag.Parse()

	if err := godotenv.Load(*envPath); err != nil {
		log.Println(".env file not loaded")
	}

	k := koanf.New(".")

	if err := k.Load(file.Provider(*configPath), yaml.Parser()); err != nil {
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

// func loadSecrets(cfg *Config) {
// 	dbUser := os.Getenv("DATABASE_USER")
// 	cfg.Database.User = dbUser

// 	dbPass := os.Getenv("DATABASE_PASSWORD")
// 	cfg.Database.Password = dbPass

// 	rmqUser := os.Getenv("RABBITMQ_USER")
// 	cfg.Mq.User = rmqUser

// 	rmqPass := os.Getenv("RABBITMQ_PASSWORD")
// 	cfg.Mq.Password = rmqPass
// }
