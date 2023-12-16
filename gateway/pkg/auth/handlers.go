package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/security"
	"gitlab.com/narm-group/service-api/api/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthHandler struct{}

var authClient authpb.UserServiceClient

func RegisterGrpcClient(url string) {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}

	authClient = authpb.NewUserServiceClient(cc)

	logrus.Infof("connection to auth service grpc: %s\n", cc.GetState().String())
}

func (h *AuthHandler) Login(c *gin.Context) {
	var creds authpb.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.String(http.StatusBadRequest, "provide valid credentials")
		return
	}

	ctx := c.Request.Context()
	res, err := authClient.Login(ctx, &creds)
	if err != nil {
		HandleError(c, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    res.Token.Value,
		Expires:  time.Unix(res.Token.ExpirationTime, 0),
		Secure:   true,
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(200, gin.H{
		"message": "login successful",
		"user_id": res.UserId,
		"token":   res.Token.Value,
		"role":    res.Role,
	})
}
func (h *AuthHandler) Signup(c *gin.Context) {
	var req authpb.SignupReq
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusUnauthorized, "invalid signup credentials")
		return
	}

	ctx := c.Request.Context()
	req.Role = 2
	res, err := authClient.Signup(ctx, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    res.Token.Value,
		Expires:  time.Unix(res.Token.ExpirationTime, 0),
		Secure:   true,
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(201, &gin.H{
		"user_id": res.UserId,
		"token":   res.Token.Value,
		"role":    res.Role,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	}
	http.SetCookie(c.Writer, cookie)
	c.String(200, "logout successful")
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	token, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "cookie is not set")
		return
	}

	ctx := c.Request.Context()
	newToken, err := authClient.RefreshToken(
		ctx,
		&authpb.RefreshTokenReq{
			Token: token,
		},
	)

	if err != nil {
		HandleError(c, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    newToken.Value,
		Expires:  time.Unix(newToken.ExpirationTime, 0),
		Secure:   true,
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
	c.String(200, "token refreshed")
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()

	var req authpb.UserInfoReq
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	// userId, err := security.GetCurrentUserId(ctx)
	// if err != nil {
	// 	c.String(http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }

	// ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
	// 	"user_id", fmt.Sprintf("%d", userId),
	// ))

	res, err := authClient.GetUserInfo(ctx, &authpb.UserInfoReq{Id: req.GetId()})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req authpb.ChangePassReq
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()
	ctx, err := addUserIdToCtx(ctx)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	_, err = authClient.ChangePassword(ctx, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *AuthHandler) EditUserProfile(c *gin.Context) {
	var req authpb.UserProfile
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()
	userId, err := security.GetCurrentUserId(ctx)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user_id", fmt.Sprintf("%d", userId),
	))

	_, err = authClient.EditUserProfile(ctx, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func HandleError(c *gin.Context, err error) {
	errStatus, ok := status.FromError(err)
	if !ok {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var status int
	switch errStatus.Code() {
	case codes.PermissionDenied:
		status = http.StatusForbidden
	case codes.Unauthenticated:
		status = http.StatusUnauthorized
	case codes.FailedPrecondition, codes.AlreadyExists, codes.InvalidArgument:
		status = http.StatusBadRequest
	case codes.NotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{
		"message": errStatus.Message(),
	})
}

func addUserIdToCtx(ctx context.Context) (context.Context, error) {
	userId, err := security.GetCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user_id", fmt.Sprintf("%d", userId),
	))

	return ctx, nil
}
