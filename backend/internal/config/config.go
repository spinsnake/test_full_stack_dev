package config

import (
	"bufio"
	"fmt"
	"net/url"
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
	MockData             bool
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
		AppPort:              getEnvAny([]string{"APP_PORT", "PORT"}, "8080"),
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
		MockData:             getEnvAnyBool([]string{"MOCKDATA", "MOCK_DATA"}, false),
		DBMaxOpenConns:       getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:       getEnvInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxLifetime:    time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		DBConnMaxIdleTime:    time.Duration(getEnvInt("DB_CONN_MAX_IDLE_TIME_MIN", 10)) * time.Minute,
	}

	dsn, err := loadDatabaseDSN()
	if err != nil {
		return Config{}, err
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

func loadDatabaseDSN() (string, error) {
	raw := strings.TrimSpace(getEnvAny([]string{"DATABASE_DSN", "DATABASE_URL", "MYSQL_URL"}, ""))
	if raw != "" {
		return normalizeDatabaseDSN(raw)
	}

	return buildMySQLDSN(), nil
}

func normalizeDatabaseDSN(raw string) (string, error) {
	if !strings.HasPrefix(raw, "mysql://") {
		return raw, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse mysql url: %w", err)
	}

	user := "root"
	password := ""
	if parsed.User != nil {
		user = parsed.User.Username()
		password, _ = parsed.User.Password()
	}

	host := parsed.Hostname()
	if host == "" {
		host = "127.0.0.1"
	}

	port := parsed.Port()
	if port == "" {
		port = "3306"
	}

	database := strings.TrimPrefix(parsed.Path, "/")
	params := parsed.Query()
	applyDefaultMySQLParams(params)

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		user,
		password,
		host,
		port,
		database,
		params.Encode(),
	), nil
}

func buildMySQLDSN() string {
	user := getEnvAny([]string{"MYSQL_USER", "MYSQLUSER"}, "root")
	password := getEnvAny([]string{"MYSQL_PASSWORD", "MYSQLPASSWORD"}, "")
	host := getEnvAny([]string{"MYSQL_HOST", "MYSQLHOST"}, "127.0.0.1")
	port := getEnvAny([]string{"MYSQL_PORT", "MYSQLPORT"}, "3306")
	database := getEnvAny([]string{"MYSQL_DATABASE", "MYSQLDATABASE"}, "gallery_db")
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

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return fallback
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

func getEnvAnyBool(keys []string, fallback bool) bool {
	for _, key := range keys {
		raw := strings.TrimSpace(os.Getenv(key))
		if raw == "" {
			continue
		}

		value, err := strconv.ParseBool(raw)
		if err != nil {
			return fallback
		}

		return value
	}

	return fallback
}

func applyDefaultMySQLParams(params url.Values) {
	if !params.Has("parseTime") {
		params.Set("parseTime", "true")
	}
	if !params.Has("loc") {
		params.Set("loc", "Local")
	}
	if !params.Has("charset") {
		params.Set("charset", "utf8mb4")
	}
}
