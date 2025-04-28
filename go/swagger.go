/*
 * Grocery List Management API
 *
 * Swagger UI configuration and setup
 */

package openapi

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui"
)

// SetupSwagger configures and adds Swagger UI routes to the router
func SetupSwagger(router *gin.Engine) {
	// Get the path to the OpenAPI YAML file
	apiSpecPath := "/api/openapi.yaml"

	// Create a new Swagger UI handler
	swaggerHandler := swgui.New(
		"Grocery List API", // Title for the Swagger UI page
		apiSpecPath,        // Path to the OpenAPI specification
		"/swagger/",        // Base URL path for Swagger UI
	)

	// Serve the OpenAPI specification file
	router.GET(apiSpecPath, func(c *gin.Context) {
		c.File(filepath.Join("api", "openapi.yaml"))
	})

	// Mount the Swagger UI at /swagger/ endpoint
	router.GET("/swagger/*any", gin.WrapH(swaggerHandler))

	// Add a redirect from /swagger to /swagger/
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
}
