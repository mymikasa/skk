package config

import "os"

type Config struct {
	DBPath string
	Port   string
}

func Load() *Config {
	dbPath := os.Getenv("SKK_DB_PATH")
	if dbPath == "" {
		dbPath = "data/skk.db"
	}
	port := os.Getenv("SKK_PORT")
	if port == "" {
		port = ":8080"
	}
	return &Config{
		DBPath: dbPath,
		Port:   port,
	}
}
