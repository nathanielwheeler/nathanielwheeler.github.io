package main

import (
	"fmt"
	"os"
	"strings"
)

// PostgresConfig holds database connection info.
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"name"`
}

// DefaultPostgresConfig will return development database information
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     checkDBEnv("host"),
		User:     checkDBEnv("user"),
		Password: checkDBEnv("password"),
		Port:     checkDBEnv("port"),
		DBName:   checkDBEnv("name"),
	}
}

// Dialect will return the dialect that GORM will use.
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionString will return a string used to connect to the database
func (c PostgresConfig) ConnectionString() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.DBName,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s, dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

// Config holds configuration variables
type Config struct {
	Port      int
	Env       string
	CSRFBytes int
}

// DefaultConfig sets up Config for a development environment
func DefaultConfig() Config {
	return Config{
		Port:      3000,
		Env:       "dev",
		CSRFBytes: 32,
	}
}

// IsProd sets Config Env to Production
func (c Config) IsProd() bool {
	return c.Env == "prod"
}

// #region DB HELPERS

func checkDBEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}

// #endregion