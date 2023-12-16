package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/security"
	"gitlab.com/narm-group/service-api/api/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func JwtValidation(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")

	if len(authorization) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authoization header is empty",
		})
		return
	}

	parts := strings.Split(authorization, " ")

	if !strings.HasPrefix(authorization, "JWT") || len(parts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authorization header doesn't contain valid jwt token",
		})
		return
	}

	ctx := c.Request.Context()
	res, err := authClient.ValidateToken(
		ctx,
		&authpb.ValidationReq{
			Token: parts[1],
		},
	)

	if err != nil {
		errStatus, _ := status.FromError(err)
		switch errStatus.Code() {
		case codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": errStatus.Message(),
			})
		case codes.Unauthenticated:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": errStatus.Message(),
			})
		default:
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	//ctx := context.WithValue(c.Request.Context(), "userId", res.UserId)
	c.Request = c.Request.WithContext(security.NewUserContext(c.Request.Context(), res.UserId))
	logrus.Infoln("setting user id : ", res.UserId)

	// c.Set("userId", res.UserId)
	// c.Set("username", res.Username)

	c.Set("role", res.Role)

	c.Next()
}

func RequireRole(role int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetInt64("role")
		if userRole != role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
