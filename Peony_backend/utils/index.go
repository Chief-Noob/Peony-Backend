package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func remote_server_load() (string, string) {
	godotenv.Load()

	host := os.Getenv("HOST_NAME")
	secret_key := os.Getenv("SECRET_KEY")
	return host, secret_key
}
