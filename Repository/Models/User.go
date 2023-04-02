package models

type User struct {
	UserId    int64
	Firstname string
	Lastname  string
	Email     string
	Links     []Link
}
