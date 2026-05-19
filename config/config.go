package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
	Supabase    Supabase
}

func Load() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	supabase := loadSupabase()

	return Config{
		Port:        port,
		JWTSecret:   jwtSecret,
		DatabaseURL: supabase.DatabaseURL(),
		Supabase:    supabase,
	}
}
