package config

import (
    "fmt"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    PublicHost             string
    Port                   string
    DBUser                 string
    DBPassword             string
    DBHost                 string
    DBPort                 string
    DBName                 string
    DBConnectionURL        string
    JWTSecret              string
    JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
    godotenv.Load()

    dbUser := getEnv("DB_USER", "pay")
    dbPassword := getEnv("DB_PASSWORD", "pay")
    dbHost := getEnv("DB_HOST", "127.0.0.1")
    dbPort := getEnv("DB_PORT", "5432")
    dbName := getEnv("DB_NAME", "app")

    dbURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        dbUser,
        dbPassword,
        dbHost,
        dbPort,
        dbName,
    )

    return Config{
        PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
        Port:       getEnv("PORT", "8080"),

        DBUser:   dbUser,
        DBPassword: dbPassword,
        DBHost:   dbHost,
        DBPort:   dbPort,
        DBName:   dbName,

        DBConnectionURL: dbURL,

        JWTSecret:              getEnv("JWT_SECRET", "secret"),
        JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
    }
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
    if value, ok := os.LookupEnv(key); ok {
        i, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return fallback
        }
        return i
    }
    return fallback
}
