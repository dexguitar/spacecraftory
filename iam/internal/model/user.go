package model

type User struct {
	UUID     string
	Info     UserInfo
	Password string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []string
}
