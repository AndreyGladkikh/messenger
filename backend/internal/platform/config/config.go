package config

import (
	"fmt"
	"log"
	"messenger/messenger/internal/platform/utils"
	"os"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Env        string
	Database   *DatabaseConfig
	Cache      *CacheConfig
	HttpServer *HttpServerConfig
}

type DatabaseConfig struct {
	Driver   string
	Schema   string
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
}

type CacheConfig struct {
	DSN string
}

type HttpServerConfig struct {
	Addr string
}

func Load() *Config {
	projRoot, err := utils.LoadProjectRoot()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load project root: %w", err))
	}

	loadEnv()

	k := koanf.New(".")

	if err := k.Load(file.Provider(projRoot+"/configs/config.yaml"), yaml.Parser()); err != nil {
		log.Fatal(err.Error())
	}

	var cfg *Config

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	loadSecrets(cfg)

	return cfg
}

func loadEnv() {
	projRoot := utils.ProjectRoot()

	env := os.Getenv("MESSENGER_ENV")
	if "" == env {
		env = "dev"
	}

	err := godotenv.Load(projRoot + "/.env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	godotenv.Load(projRoot + "/.env." + env)

	if "test" != env {
		godotenv.Load(projRoot + "/.env.local")
	}

	godotenv.Load(projRoot + "./env." + env + ".local")
}

func loadSecrets(cfg *Config) {
	dbUser := os.Getenv("DATABASE_USER")
	cfg.Database.User = dbUser

	dbPass := os.Getenv("DATABASE_PASSWORD")
	cfg.Database.Password = dbPass
}
