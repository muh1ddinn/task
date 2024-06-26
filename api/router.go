package api

import (
	"errors"
	"net/http"
	"task/api/handler"
	"task/api/middleware"
	"task/pkg/logger"
	"task/service"

	"github.com/gin-gonic/gin"

	_ "task/api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceMangaer, log logger.ILogger) *gin.Engine {
	h := handler.NewStrg(services, log)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(authMiddleware)
	r.Use(middleware.RateLimiter(1, 5))

	r.POST("/contact", h.Create)
	r.GET("/contact/:id", h.GetByID)
	r.GET("/contact", h.Getall)
	r.DELETE("/contact/:id", h.Deletehard)
	r.PATCH("/contact", h.Patch)
	r.DELETE("contact_s/:id", h.Deletesoft)

	r.GET("/contact/history/:id", h.History)

	r.GET("/contacts/export/csv", h.ExportToCSV)
	r.POST("/contacts/import", h.ImportContacts)

	r.POST("/categories", h.Createcat)
	r.GET("/categories/:id", h.GetByIDcat)
	r.GET("/categories", h.Getallcat)
	r.DELETE("/categories/:id", h.Deletehardcat)
	r.PATCH("/categories", h.Patchcat)
	r.DELETE("categories_s/:id", h.Deletesoftcat)

	return r
}

func authMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized, you can use any character for ApiKeyAuth (apiKey)"))
	}
	c.Next()
}
