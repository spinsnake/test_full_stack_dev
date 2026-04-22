package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName              string
	AppHost              string
	AppPort              string
	APIPrefix            string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	CORSAllowOrigins     string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool
	DefaultPageLimit     int
	MaxPageLimit         int
	DatabaseDSN          string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetime    time.Duration
	DBConnMaxIdleTime    time.Duration
}

func Load() (Config, error) {
	loadEnvFiles(".env", filepath.Join("backend", ".env"))

	cfg := Config{
		AppName:              getEnv("APP_NAME", "image-gallery-api"),
		AppHost:              getEnv("APP_HOST", "0.0.0.0"),
		AppPort:              getEnv("APP_PORT", "8080"),
		APIPrefix:            getEnv("API_PREFIX", "/api"),
		ReadTimeout:          time.Duration(getEnvInt("READ_TIMEOUT_SEC", 15)) * time.Second,
		WriteTimeout:         time.Duration(getEnvInt("WRITE_TIMEOUT_SEC", 15)) * time.Second,
		IdleTimeout:          time.Duration(getEnvInt("IDLE_TIMEOUT_SEC", 60)) * time.Second,
		CORSAllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		CORSAllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", false),
		DefaultPageLimit:     getEnvInt("DEFAULT_PAGE_LIMIT", 12),
		MaxPageLimit:         getEnvInt("MAX_PAGE_LIMIT", 60),
		DBMaxOpenConns:       getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:       getEnvInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxLifetime:    time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		DBConnMaxIdleTime:    time.Duration(getEnvInt("DB_CONN_MAX_IDLE_TIME_MIN", 10)) * time.Minute,
	}

	dsn := strings.TrimSpace(os.Getenv("DATABASE_DSN"))
	if dsn == "" {
		dsn = buildMySQLDSN()
	}
	if dsn == "" {
		return Config{}, fmt.Errorf("database dsn is empty")
	}
	cfg.DatabaseDSN = dsn

	return cfg, nil
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%s", c.AppHost, c.AppPort)
}

func buildMySQLDSN() string {
	user := getEnv("MYSQL_USER", "root")
	password := getEnv("MYSQL_PASSWORD", "")
	host := getEnv("MYSQL_HOST", "127.0.0.1")
	port := getEnv("MYSQL_PORT", "3306")
	database := getEnv("MYSQL_DATABASE", "gallery_db")
	params := strings.TrimPrefix(getEnv("MYSQL_PARAMS", "parseTime=true&loc=Local&charset=utf8mb4"), "?")

	if password == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, database, params)
}

func loadEnvFiles(paths ...string) {
	for _, path := range paths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err != nil {
			continue
		}
		loadEnvFile(path)
	}
}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		_ = os.Setenv(key, value)
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func getEnvBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}
