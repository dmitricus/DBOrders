package model

import "time"

// User is ...
type User struct {
	ID       int64
	Username string
	Password string
	Created  time.Time
	Email    string
	IsAdmin  bool
	Title    string
}
