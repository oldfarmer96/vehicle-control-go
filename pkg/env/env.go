// Package env - variables de entorno
package env

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AppEnv      string
	CORSURLs    string
	CookieName  string
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", errors.New("la variable " + key + " es requerida")
	}
	return value, nil
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error al cargar el archivo .env")
	}

	port, err := getEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	db, err := getEnv("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	appEnv, err := getEnv("APP_ENV")
	if err != nil {
		return Config{}, err
	}

	cookieName, err := getEnv("COOKIE_NAME")
	if err != nil {
		return Config{}, err
	}

	corsURLs, _ := os.LookupEnv("CORS_URLS")

	return Config{
		Port:        port,
		DatabaseURL: db,
		AppEnv:      appEnv,
		CORSURLs:    corsURLs,
		CookieName:  cookieName,
	}, nil
}
