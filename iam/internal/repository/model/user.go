package model

// UserRow is a flat struct that maps directly to the database columns
type UserRow struct {
	ID                  string   `db:"id"`
	Login               string   `db:"login"`
	Email               string   `db:"email"`
	NotificationMethods []string `db:"notification_methods"`
	Password            string   `db:"password"`
}

// User is the domain model with nested structure
type User struct {
	UUID     string
	Password string
	Info     UserInfo
}

// UserInfo is a nested struct that contains user information
type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []string
}
