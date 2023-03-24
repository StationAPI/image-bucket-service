package routes

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
	"storj.io/uplink"
)

type requestImage struct {
	Data string `json:"image"`
}

func Create(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	requestImage := requestImage{}

	err := ProcessBody(r.Body, &requestImage)

	if err != nil {
		http.Error(w, "there was an error processing the request body", http.StatusInternalServerError)

		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(requestImage.Data)

	if err != nil {
		http.Error(w, "there was an error processing the image", http.StatusInternalServerError)

		return err
	}

	reader := bytes.NewReader(decoded)

	var image image.Image
	buf := new(bytes.Buffer)

	switch strings.TrimSuffix(requestImage.Data[5:strings.Index(requestImage.Data, ",")], ";base64") {
	case "image/png":
		image, err = png.Decode(reader)

		if err != nil {
			http.Error(w, "there was an error processing the image", http.StatusInternalServerError)

			return err
		}

		err := png.Encode(buf, image)

		if err != nil {
			http.Error(w, "there was an error processing the image", http.StatusInternalServerError)
		}
	case "image/jpeg":
		image, err = jpeg.Decode(reader)

		if err != nil {
			http.Error(w, "there was an error processing the image", http.StatusInternalServerError)

			return err
		}

		err := jpeg.Encode(buf, image, nil)

		if err != nil {
			http.Error(w, "there was an error processing the image", http.StatusInternalServerError)
		}
	}

	ctx := context.Background()

	access, err := uplink.ParseAccess(os.Getenv("STORJ_TOKEN"))

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	project, err := uplink.OpenProject(ctx, access)

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	_, err = project.EnsureBucket(ctx, "station-dev")

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	upload, err := project.UploadObject(ctx, "station-dev", os.Getenv("STORJ_UPLOAD_KEY"), nil)

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	_, err = io.Copy(upload, buf)

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	err = upload.Commit()

	if err != nil {
		http.Error(w, "there was an error uploading the image", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write([]byte("the image has been successfully updated"))

	return nil
}
