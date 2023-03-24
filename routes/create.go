package routes

import (
	"net/http"

	"gorm.io/gorm"
  neon "github.com/stationapi/image-bucket-service/db"
)

func Create(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
  image := neon.Image{} 

  err := ProcessBody(r.Body, &image)

  if err != nil {
    http.Error(w, "there was an error processing the request body", http.StatusInternalServerError)

    return err
  }

  neon.CreateImage(image, db)

  w.WriteHeader(200)
  w.Write([]byte("the image has been successfully updated"))

  return nil
}
