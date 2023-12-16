package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/narm-group/gateway/internal/server"
	"gitlab.com/narm-group/gateway/pkg/auth"
	"gitlab.com/narm-group/gateway/pkg/book"
	"gitlab.com/narm-group/gateway/pkg/offers"
	"gitlab.com/narm-group/gateway/pkg/review"
)

type Role int64

const (
	Admin Role = iota + 1
	User
)

func InitRoutes(s *server.Server) {
	public := s.Router.Group("/api/v1")
	protected := s.Router.Group("/api/v1")

	//handlers
	offersHandler := &offers.OfferHandler{}
	reviewHandler := &review.ReviewHandler{}
	authHandler := &auth.AuthHandler{}
	bookHandler := &book.BookHandler{}

	//jwt middleware
	protected.Use(auth.JwtValidation)

	//public routes
	authG := public.Group("/auth")
	authG.POST("/login", authHandler.Login)
	authG.POST("/signup", authHandler.Signup)
	authG.GET("/logout", authHandler.Logout)
	authG.POST("/refresh", authHandler.Refresh)

	//user related services
	userG := protected.Group("/users")
	userG.POST("/update", authHandler.EditUserProfile)
	userG.POST("/change-password", authHandler.ChangePassword)

	public.POST("/user-info", authHandler.GetUserInfo)
	//	userG.GET("/info", authHandler.GetUserInfo)

	//protected routes
	ag := protected.Group("/offers")
	ag.POST("", offersHandler.SubmitOffer)

	protected.POST("/pending-offers", auth.RequireRole(int64(Admin)), reviewHandler.GetPendingOffers)
	protected.POST("/user-offers", reviewHandler.GetUserOffers)

	rg := protected.Group("/review")
	rg.POST("", auth.RequireRole(int64(Admin)), reviewHandler.SubmitReviewResult)

	bg := public.Group("/books")
	bg.POST("", bookHandler.GetBookOffers)
	bg.POST("/info", bookHandler.GetBookInfo)

	protected.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello baby, it is uppppppppppppppp")
	})
}
