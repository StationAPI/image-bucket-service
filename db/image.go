package db

import "gorm.io/gorm"

type Image struct {
	Data string `json:"data"` // b64 encoded image data
}

func CreateImage(image Image, db gorm.DB) {
	db.Create(image)
}
