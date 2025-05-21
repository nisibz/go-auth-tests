package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Container contains environment variables for the application, and http server
type (
	Container struct {
		App          *App
		HTTP         *HTTP
		Mongo        *Mongo
		JwtSecretKey *JWT
	}

	// App contains all the environment variables for the application
	App struct {
		Name string
		Env  string
	}

	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}

	// Mongo contains all the environment variables for MongoDB
	Mongo struct {
		URI     string
		DB_NAME string
	}

	JWT struct {
		JWT_SECRET_KEY string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	mongo := &Mongo{
		URI:     os.Getenv("DB_URI"),
		DB_NAME: os.Getenv("DB_NAME"),
	}

	jwt := &JWT{
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
	}

	return &Container{
		app,
		http,
		mongo,
		jwt,
	}, nil
}
