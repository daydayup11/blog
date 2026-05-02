package config

import "os"

type Config struct {
	DBPath      string
	JWTSecret   string
	Port        string
	AdminUser   string
	AdminPass   string
	FrontendDir string
}

func Load() Config {
	return Config{
		DBPath:      getEnv("DB_PATH", "./data/blog.db"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-prod"),
		Port:        getEnv("PORT", "8080"),
		AdminUser:   getEnv("ADMIN_USER", "admin"),
		AdminPass:   getEnv("ADMIN_PASS", "admin123"),
		FrontendDir: getEnv("FRONTEND_DIR", "../frontend"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
