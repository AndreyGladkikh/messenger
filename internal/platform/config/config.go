package config

type Config struct {
	Env   string
	DB    DatabaseConfig
	Cache DatabaseConfig
}

type DatabaseConfig struct {
	Driver string
	DSN string
}

type CacheConfig struct {
	DSN string
}

func Load() *Config {
	return &Config{}
}
