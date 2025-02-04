package entity

import "time"

// тут тоже надо поправить, когда store будут реализовывать
type Store struct {
	ID          uint
	Name        string
	Description string
	Products    []Product
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
