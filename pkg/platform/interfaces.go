package platform

import "github.com/gin-gonic/gin"

type Platform interface {
	Name() string
	// Sync synchronizes the platform with the latest data.
	// It should be called periodically to keep the platform up-to-date.
	Sync() error

	// RegisterRoutes registers the routes for the platform.
	RegisterRoutes(gin.IRouter)

	// Close closes the platform and releases any resources.
	Close() error
}
