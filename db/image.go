package db

type Image struct {
	Data string `json:"image"` // b64 encoded image data
}

func CreateImage(image Image, db gorm.DB) {
	db.Create(image)
}