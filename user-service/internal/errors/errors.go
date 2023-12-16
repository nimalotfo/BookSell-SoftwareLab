package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUsernameAlreadyExists     = status.Errorf(codes.AlreadyExists, "username already exists")
	ErrMobileNumberAlreadyExists = status.Errorf(codes.AlreadyExists, "mobile number already exists")
	ErrUserNotFound              = status.Errorf(codes.NotFound, "user not found")
	ErrWrongPassword             = status.Errorf(codes.PermissionDenied, "password is incorrect")
	ErrInternal                  = status.Errorf(codes.Internal, "internal error")
	ErrEmptyCredentials          = status.Errorf(codes.InvalidArgument, "credentials must not be empty")
	ErrFetchingUserInfo          = status.Errorf(codes.Internal, "error fetching user info")
	ErrFetchUserContext          = status.Errorf(codes.Internal, "couldn't extract userId from context")
	ErrUsernameNotFound          = status.Errorf(codes.NotFound, "user with this username not found")
	ErrUserPermissionDenied      = status.Errorf(codes.PermissionDenied, "permission denied to perform this operation")
	ErrHashingPassword           = status.Errorf(codes.Internal, "error processing credentials")
	ErrUpdatingUserInfo          = status.Errorf(codes.Internal, "error updating user info")
)
