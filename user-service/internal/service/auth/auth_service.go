package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	authpb "gitlab.com/narm-group/service-api/api/authpb"
	"gitlab.com/narm-group/user-service/internal/config"
	"gitlab.com/narm-group/user-service/internal/database"
	"gitlab.com/narm-group/user-service/internal/encryption"
	api_errors "gitlab.com/narm-group/user-service/internal/errors"
	"gitlab.com/narm-group/user-service/internal/models"
	"gitlab.com/narm-group/user-service/internal/persistence"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	authpb.UnimplementedUserServiceServer
}

func RegisterGrpcService(s *grpc.Server) {
	authpb.RegisterUserServiceServer(s, &AuthService{})
}

type Claims struct {
	UserId      int64               `json:"user_id"`
	Username    string              `json:"username"`
	Roles       []models.Role       `json:"roles"`
	Permissions []models.Permission `json:"permissions"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *AuthService) Login(ctx context.Context, creds *authpb.Credentials) (*authpb.AuthRes, error) {
	if len(creds.Username) == 0 || len(creds.Password) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "credentials must not be empty")
	}

	user, err := persistence.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		if errors.Is(err, api_errors.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	if !encryption.CheckPasswordHash(creds.Password, user.PasswordHash) {
		return nil, status.Errorf(codes.Unauthenticated, "wrong credentials")
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	jwtKey := config.GetCfg().JwtKey

	err = loadRoles(ctx, user)
	if err != nil {
		return nil, err
	}

	tokenStr, err := generateToken(expiresAt, user, jwtKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating token")
	}

	token := &authpb.Token{
		Value:          tokenStr,
		ExpirationTime: expiresAt.Unix(),
	}

	return &authpb.AuthRes{
		UserId: user.ID,
		Token:  token,
		Role:   user.Roles[0].ID,
	}, nil
}

func (s *AuthService) Signup(ctx context.Context, req *authpb.SignupReq) (*authpb.AuthRes, error) {
	if len(req.Username) == 0 || len(req.Password) == 0 || len(req.MobileNumber) == 0 {
		return nil, api_errors.ErrEmptyCredentials
	}

	hashedPass, err := encryption.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing passwrod")
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPass,
		MobileNumber: req.MobileNumber,
	}

	insertedId, err := persistence.CreateUser(ctx, user, req.Role)
	if err != nil {
		fmt.Println("here")
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		fmt.Println("here2")
		// if errors.Is(err, api_errors.ErrUsernameAlreadyExists) {
		// 	return nil, status.Errorf(codes.AlreadyExists, "this username already exists")
		// }
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "error creating new user")
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	jwtKey := config.GetCfg().JwtKey
	err = loadRoles(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Printf("userrrrr -> %#v\n", user)

	tokenStr, err := generateToken(expiresAt, user, jwtKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating token")
	}

	return &authpb.AuthRes{
		UserId: insertedId,
		Token: &authpb.Token{
			Value:          tokenStr,
			ExpirationTime: expiresAt.Unix(),
		},
		Role: user.Roles[0].ID,
	}, nil

}

func loadRoles(ctx context.Context, user *models.User) error {
	db := database.GetDB(ctx)
	// var permissions []string

	err := db.Model(&models.User{}).
		Preload("Roles").
		Preload("Roles.Permissions").
		Where("id = ?", user.ID).
		Take(&user).Error

	return err

	// err := db.Table("users AS u").
	// 	Joins("JOIN user_roles AS ur ON ur.user_id = u.id").
	// 	Joins("JOIN role_permissions AS rp ON rp.role_id = ur.role_id").
	// 	Joins("JOIN permissions AS p ON p.id = rp.permission_id").
	// 	Where("u.id = ?", user.ID).
	// 	Select("p.title").
	// 	Take(&permissions).
	// 	Error

	// return permissions, err
}

func (s *AuthService) ValidateToken(ctx context.Context, req *authpb.ValidationReq) (*authpb.ValidationRes, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		req.Token,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.GetCfg().JwtKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
		}
		logrus.Errorf("error validating token: %v\n", err)
		return nil, status.Errorf(codes.FailedPrecondition, "token is invalid")
	}

	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	return &authpb.ValidationRes{
		UserId:   claims.UserId,
		Username: claims.Username,
		Role:     claims.Roles[0].ID,
		//	Permissions: claims.Permissions,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *authpb.RefreshTokenReq) (*authpb.Token, error) {
	jwtKey := config.GetCfg().JwtKey

	fmt.Println("received token is", req.Token)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		req.Token,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
		}
		logrus.Errorf("error validating token: %v\n", err)
		return nil, status.Errorf(codes.FailedPrecondition, "token is invalid")
	}

	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	if time.Until(claims.ExpiresAt.Time) > 2*time.Minute {
		return nil, status.Errorf(codes.FailedPrecondition, "token can't get refreshed now")
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := newToken.SignedString([]byte(jwtKey))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error signing token")
	}

	return &authpb.Token{
		Value:          tokenStr,
		ExpirationTime: expiresAt.Unix(),
	}, nil

}

func generateToken(expiresAt time.Time, user *models.User, jwtKey string) (tokenStr string, err error) {
	permissions := []models.Permission{}
	for _, r := range user.Roles {
		permissions = append(permissions, r.Permissions...)
	}
	claims := &Claims{
		UserId:      user.ID,
		Username:    user.Username,
		Roles:       user.Roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(jwtKey))
	return
}

func (s *AuthService) GetUserInfo(ctx context.Context, req *authpb.UserInfoReq) (*authpb.UserInfoRes, error) {
	user, err := persistence.GetUserById(ctx, req.Id)
	if err != nil {
		logrus.Errorf("error while getting user by id: %v\n", err)
		return nil, api_errors.ErrUsernameNotFound
	}

	return &authpb.UserInfoRes{
		Id:           user.ID,
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		Email:        user.Email,
		City:         user.City,
	}, nil
}

func (s *AuthService) EditUserProfile(ctx context.Context, req *authpb.UserProfile) (res *emptypb.Empty, err error) {
	res = &emptypb.Empty{}
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		err = api_errors.ErrUserPermissionDenied
		return
	}

	updatedUser := models.User{
		Username:     req.Username,
		Email:        req.Email,
		City:         req.City,
		MobileNumber: req.MobileNumber,
	}

	db := database.GetDB(ctx)
	err = db.Model(&models.User{}).
		Where("id = ?", userId).
		Updates(&updatedUser).
		Error

	if err != nil {
		logrus.Error("error updating user info: %v\n", err)
		err = api_errors.ErrUpdatingUserInfo
		return
	}

	return
}

func (s *AuthService) ChangePassword(ctx context.Context, req *authpb.ChangePassReq) (res *emptypb.Empty, err error) {
	res = &emptypb.Empty{}
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		err = api_errors.ErrUserPermissionDenied
		return
	}

	err = persistence.ChangePassword(ctx, userId, req.PrevPassword, req.NewPassword)
	if err != nil {
		logrus.Errorf("error changing password: %v\n", err)
		return nil, err
	}

	return
}

func getUserIdFromCtx(ctx context.Context) (int64, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	userIdStr := md["user_id"][0]

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logrus.Errorf("user id is invalid %s", userIdStr)
		return 0, api_errors.ErrUserPermissionDenied

	}
	return int64(userId), nil
}
