package http

import (
	"log/slog"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nisibz/go-auth-tests/internal/adapter/config"
	sloggin "github.com/samber/slog-gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	config *config.HTTP,
) (*Router, error) {
	if config.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "welcome to go-auth-tests"})
	})

	return &Router{router}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
