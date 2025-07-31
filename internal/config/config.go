package config

import (
	"cmp"
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host        string
	Port        int
	DbDSN       string
	MigratePath string
	Debug       bool
}

const (
	defaultHost = "0.0.0.0"
	defaultPort = 8081
)

func ReadConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", defaultHost, "server host address")
	flag.IntVar(&cfg.Port, "port", defaultPort, "server port")
	flag.BoolVar(&cfg.Debug, "debug", false, "enable logger debug level")
	flag.Parse()

	cfg.Host = cmp.Or(os.Getenv("SRV_HOST"), cfg.Host)
	if tmp := os.Getenv("SRV_PORT"); tmp != "" {
		port, err := strconv.Atoi(tmp)
		if err != nil {
			log.Println(err.Error())
			return cfg
		}
		cfg.Port = port
	}
	cfg.MigratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), "migrations")
	cfg.DbDSN = cmp.Or(os.Getenv("DB_DSN"), "postgres://user:password@localhost:5432/test?sslmode=disable")
	return cfg
}
