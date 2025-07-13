package config

import "os"

// Config holds application configuration loaded from environment variables.
type Config struct {
    AppEnv     string
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    SentryDSN  string
}

func Load() (*Config, error) {
    cfg := &Config{
        AppEnv:     getenv("APP_ENV", "dev"),
        DBHost:     getenv("DB_HOST", "localhost"),
        DBPort:     getenv("DB_PORT", "5432"),
        DBUser:     getenv("DB_USER", "postgres"),
        DBPassword: getenv("DB_PASSWORD", "postgres"),
        DBName:     getenv("DB_NAME", "filmhub"),
        JWTSecret:  getenv("JWT_SECRET", "supersecretkey"),
        SentryDSN:  getenv("SENTRY_DSN", ""),
    }
    return cfg, nil
}

func getenv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
} 