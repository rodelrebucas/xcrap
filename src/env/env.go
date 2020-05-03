package env

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

// Var defines environment variables
type Var struct {
	Port, Env, Secret, RedisHost, RedisPass string
}

// LoadEnv loads all environment variables in .env file
func LoadEnv() *Var {
	// Load env variables when in development
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}

	return &Var{
		Port:   os.Getenv("PORT"),
		Env:     os.Getenv("ENV"),
		Secret:  os.Getenv("SECRET_KEY"),
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPass: os.Getenv("REDIS_PASS"),
	}
}
