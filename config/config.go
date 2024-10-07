package config

import (
	"os"
)

type ServerConfig struct {
	Address string
}

type LogConfig struct {
	DebugMode bool
}

type PgConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewRemoteServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: getEnv("REMOTE_SERVER_ADDRESS", ""),
	}
}

func NewInternalServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: getEnv("INTERNAL_SERVER_ADDRESS", ""),
	}
}

func NewLogMode() *LogConfig {
	return &LogConfig{
		DebugMode: getEnvAsBool("SERVER_DEBUG_MODE", true),
	}
}

func NewRemotePgConfig() *PgConfig {
	return &PgConfig{
		Host:     getEnv("POSTGRES_REMOTE_SERVER_HOST", ""),
		Port:     getEnv("POSTGRES_REMOTE_SERVER_PORT", ""),
		Username: getEnv("POSTGRES_REMOTE_SERVER_USERNAME", ""),
		Password: getEnv("POSTGRES_REMOTE_SERVER_PASSWORD", ""),
		Database: getEnv("POSTGRES_REMOTE_SERVER_DB_NAME", ""),
	}
}

func NewInternalPgConfig() *PgConfig {
	return &PgConfig{
		Host:     getEnv("POSTGRES_INTERNAL_SERVER_HOST", ""),
		Port:     getEnv("POSTGRES_INTERNAL_SERVER_PORT", ""),
		Username: getEnv("POSTGRES_INTERNAL_SERVER_USERNAME", ""),
		Password: getEnv("POSTGRES_INTERNAL_SERVER_PASSWORD", ""),
		Database: getEnv("POSTGRES_INTERNAL_SERVER_DB_NAME", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "false" {
			return false
		}
		if value == "true" {
			return true
		}
	}
	return defaultVal
}
