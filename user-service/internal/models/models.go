package models

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	MobileNumber string
	Email        string
	City         string
	Roles        []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID          int64
	Title       string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID    int64
	Title string
}

// const (
// 	ADMIN Role = iota + 1
// 	USER
// 	REVIEWER
// )
