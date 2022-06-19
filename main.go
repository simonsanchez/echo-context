package main

import (
	"log"
	"os"

	"github.com/simonsanchez/echo-context/api"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// ***************** ENV VARS *****************
	port := getEnvVar("PORT", "8080")

	server := api.Server{}

	return server.Listen(port)
}

func getEnvVar(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	if fallback == "" {
		panic("environment variable " + key + " must be set")
	}

	return fallback
}
