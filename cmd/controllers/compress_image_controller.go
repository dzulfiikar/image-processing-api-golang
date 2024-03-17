package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/services"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/validations"
	"github.com/gin-gonic/gin"
)

/**
* Input:
* images: []multipart.FileHeader
* qualities: []Int
 */
func CompressImage(c *gin.Context) {
	// validate input
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return
	}

	images, err := validations.ValidateImages(form, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	imageQualities, err := validateQuality(images, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var CompressImageInputDTO dtos.CompressImageInputDTO
	
	for i, image := range images {
		CompressImageInputDTO.CompressImages = append(CompressImageInputDTO.CompressImages, &dtos.CompressImage{
			Image:   image,
			Quality: int(imageQualities[i]),
		})
	}

	result, err := services.CompressImageService(CompressImageInputDTO)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	defer os.Remove(result.FilePath)

	c.FileAttachment(result.FilePath, result.FileName)
}

func validateQuality(images []*multipart.FileHeader, c *gin.Context) ([]int, error) {
	qualities := c.PostFormArray("qualities")
	
	if len(qualities) == 0 {
		return nil, errors.New("qualities is required")
	}

	if len(images) != len(qualities) {
		return nil, errors.New("mismatch images with qualities")
	}

	parsedQualities := make([]int, 0)

	for i, quality := range qualities {
		parsedQuality, err := strconv.ParseInt(quality, 10, 64)
		if err != nil {
			return nil, errors.New("qualities must be integer")
		}
		if images[i].Header.Get("Content-Type") == "image/png" {
			if parsedQuality > 9 || parsedQuality < 0 {
				return nil, errors.New("qualities for png file(s) must be between 0 - 9")
			}
		} else {
			if parsedQuality < 0 || parsedQuality > 100 {
				return nil, errors.New("qualities for jpg file(s) must be between 0 - 100")
			}
		}

		parsedQualities = append(parsedQualities, int(parsedQuality))
	}
	return parsedQualities, nil
}
