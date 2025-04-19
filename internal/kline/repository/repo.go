package repository

import "github.com/CrazyThursdayV50/pkgo/store/db/gorm"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
