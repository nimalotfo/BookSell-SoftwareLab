package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(g *gin.RouterGroup, path string) {
	r := g.Group(path)

	authHandler := &AuthHandler{}

	r.POST("/login", authHandler.Login)
	r.POST("/signup", authHandler.Signup)
	r.GET("/logout", authHandler.Logout)
	r.POST("/refresh", authHandler.Refresh)
}
