package config

import (
	"log"
	"messenger/messenger/internal/platform/utils"
	"os"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Env      string
	Database DatabaseConfig
	Cache    CacheConfig
	HttpServer    HttpServerConfig
}

type DatabaseConfig struct {
	Driver string
	Schema string
	Host string
	Port string
	Name string
	User string
	Pass string
	DSN    string
}

type CacheConfig struct {
	DSN string
}

type HttpServerConfig struct {
	Addr string
	Port int
}

func Load() *Config {
	loadEnv()

    k := koanf.New(".")

	projRoot, _ := utils.ProjectRoot()
    if err := k.Load(file.Provider(projRoot + "/configs/config.yaml"), yaml.Parser()); err != nil {
        log.Fatal(err.Error())
    }

    var cfg *Config

    if err := k.Unmarshal("", &cfg); err != nil {
        log.Fatal(err.Error())
    }

    return cfg
	
	// s3Bucket := os.Getenv("S3_BUCKET")
	// secretKey := os.Getenv("SECRET_KEY")

	// return &Config{}
}

func loadEnv() {
	projRoot, _ := utils.ProjectRoot()

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
