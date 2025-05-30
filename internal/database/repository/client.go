package repository

import (
	"github.com/google/uuid"
	"github.com/welovemedia/ffmate/internal/database/model"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func (c *Client) GetOrCreateClient() (*model.Client, error) {
	var client *model.Client
	err := c.DB.First(&client).Error
	if err == nil {
		return client, nil
	}

	client = &model.Client{
		Uuid: uuid.NewString(),
	}

	db := c.DB.Create(client)
	return client, db.Error
}

func (c *Client) Setup() {
	c.DB.AutoMigrate(&model.Client{})
}

func (Client) TableName() string {
	return "client"
}
