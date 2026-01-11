package models

type ErrorPageData struct {
	CurrentUser *User
	Status      int
	Message     string
}
