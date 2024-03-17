package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/services"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/validations"
	"github.com/gin-gonic/gin"
)

func ResizeImage(c *gin.Context) {
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

	imageDimensions, err := validateDimensions(len(images), c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var resizeImageInputDto dtos.ResizeImageInputDTO

	for i, image := range images {
		resizeImageInputDto.Images = append(resizeImageInputDto.Images, &dtos.ResizeImage{
			Image:  image,
			Width:  imageDimensions[i].Width,
			Height: imageDimensions[i].Height,
		})
	}

	result, err := services.ResizeImageService(resizeImageInputDto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	defer os.Remove(result.FilePath)

	c.FileAttachment(result.FilePath, result.FileName)
}

func validateDimensions(imageCount int, c *gin.Context) ([]dtos.ImageDimension, error) {
	imageDimensions, success := c.GetPostFormArray("image_dimensions")
	if !success {
		return nil, errors.New("image_dimensions is required")
	}

	if imageCount != len(imageDimensions) {
		return nil, errors.New("mismatch images with image_dimensions request")
	}

	parsedDimensions := make([]dtos.ImageDimension, 0)

	for _, imageDimension := range imageDimensions {
		var parsedDimension dtos.ImageDimension
		err := json.Unmarshal([]byte(imageDimension), &parsedDimension)
		if err != nil {
			return nil, errors.New("invalid image_dimensions request")
		}

		parsedDimensions = append(parsedDimensions, parsedDimension)
	}

	return parsedDimensions, nil
}
