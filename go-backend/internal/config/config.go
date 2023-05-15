package config

import (
	"log"
)

type AppConfig struct {
	InProduction bool
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
}
