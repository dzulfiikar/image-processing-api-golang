package validations

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func ValidateImages(form *multipart.Form, c *gin.Context) ([]*multipart.FileHeader, error) {
	images := form.File["images"]
	if images == nil {

		return nil, errors.New("image(s) are required")
	}

	for i, image := range images {
		if image.Header.Get("Content-Type") != "image/png" && image.Header.Get("Content-Type") != "image/jpeg" {
			return nil, errors.New("image type must be image/png or image/jpeg")
		}

		images[i].Filename = "image_" + fmt.Sprintf("%d", i)
	}

	return images, nil
}
