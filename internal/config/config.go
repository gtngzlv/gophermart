package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var (
	RunAddress           = "RUN_ADDRESS"
	DatabaseAddress      = "DATABASE"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

type AppConfig struct {
	RunAddress           string
	DatabaseAddress      string
	AccrualSystemAddress string
}

func LoadConfig() *AppConfig {
	var err error
	config := &AppConfig{}
	getArgs(config)
	getENVs(config)
	if config.DatabaseAddress == "" {
		config.DatabaseAddress, err = returnDefaultDB()
		if err != nil {
			log.Fatal("Failed to load default DB connection")
		}
	}
	return config
}

func getArgs(cfg *AppConfig) {
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "Application run address")
	flag.StringVar(&cfg.DatabaseAddress, "d", "", "Database address")
	flag.StringVar(&cfg.DatabaseAddress, "r", "", "Accrual system address")
	flag.Parse()
}

func getENVs(cfg *AppConfig) {
	envRunAddr := strings.TrimSpace(os.Getenv(RunAddress))
	if envRunAddr != "" {
		cfg.RunAddress = envRunAddr
	}

	databaseAddr := strings.TrimSpace(os.Getenv(DatabaseAddress))
	if databaseAddr != "" {
		cfg.DatabaseAddress = databaseAddr
	}

	accrualAddr := strings.TrimSpace(os.Getenv(AccrualSystemAddress))
	if accrualAddr != "" {
		cfg.AccrualSystemAddress = AccrualSystemAddress
	}
}

func returnDefaultDB() (string, error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}
	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.username"),
		viper.GetString("database.name"),
		viper.GetString("database.password"),
		viper.GetString("database.sslmode"))
	return connection, nil
}
