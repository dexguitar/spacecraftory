package model

// UserRow is a flat struct that maps directly to the database columns
type UserRow struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

// NotificationMethodRow maps to the notification_methods table
type NotificationMethodRow struct {
	ID           string `db:"id"`
	UserUUID     string `db:"user_uuid"`
	ProviderName string `db:"provider_name"`
	Target       string `db:"target"`
}
