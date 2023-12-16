package persistence

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/user-service/internal/database"
	"gitlab.com/narm-group/user-service/internal/encryption"
	api_errors "gitlab.com/narm-group/user-service/internal/errors"
	"gitlab.com/narm-group/user-service/internal/models"

	"gorm.io/gorm"
)

func getUser(ctx context.Context, conds models.User) (*models.User, error) {
	db := database.GetDB(ctx)

	var user models.User
	err := db.Model(&models.User{}).Preload("Roles").Take(&user, conds).Error
	fmt.Printf("user : %#v\n", user)
	if err != nil {
		fmt.Printf("err -> %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Yessss")
		}
		return nil, err
	}
	return &user, err
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := getUser(ctx, models.User{Username: username})

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, api_errors.ErrUserNotFound
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func GetUserById(ctx context.Context, id int64) (*models.User, error) {
	user, err := getUser(ctx, models.User{ID: id})

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, api_errors.ErrUserNotFound
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func UserExists(ctx context.Context, user models.User) (bool, error) {
	_, err := getUser(ctx, user)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CreateUser(ctx context.Context, user *models.User, role int64) (int64, error) {
	exists, err := UserExists(ctx, models.User{Username: user.Username})
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, api_errors.ErrUsernameAlreadyExists
	}

	exists, err = UserExists(ctx, models.User{MobileNumber: user.MobileNumber})
	if err != nil {
		return 0, api_errors.ErrEmptyCredentials
	}

	if exists {
		return 0, api_errors.ErrMobileNumberAlreadyExists
	}

	db := database.GetDB(ctx)

	err = db.Table("users").Create(&user).Error
	if err != nil {
		return 0, err
	}

	err = db.Table("user_roles").
		Create(map[string]interface{}{
			"user_id": user.ID,
			"role_id": role,
		}).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func ChangePassword(ctx context.Context, userId int64, prevPass, newPass string) error {
	db := database.GetDB(ctx)

	var curPassHash string
	err := db.Model(&models.User{}).
		Where("id = ?", userId).
		Select("password_hash").
		Take(&curPassHash).
		Error

	if !encryption.CheckPasswordHash(prevPass, curPassHash) {
		return api_errors.ErrWrongPassword
	}

	newPassHash, err := encryption.HashPassword(newPass)
	if err != nil {
		logrus.Errorf("error hashing password: %v\n", err)
		return api_errors.ErrHashingPassword
	}

	err = db.Model(&models.User{}).
		Where("id = ?", userId).
		Update("password_hash", newPassHash).Error

	if err != nil {
		logrus.Errorf("error updating password: %v\n", err)
		return api_errors.ErrInternal
	}

	return nil
}
