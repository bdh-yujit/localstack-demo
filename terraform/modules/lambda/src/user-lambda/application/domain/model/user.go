package model

import "time"

type User struct {
	ID        string
	Name      string
	BirthDate time.Time
}
