package model

type Client struct {
	ID uint `gorm:"primarykey"`

	Uuid string
}
