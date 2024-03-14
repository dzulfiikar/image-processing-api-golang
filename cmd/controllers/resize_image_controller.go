package controllers

import (
	"net/http"
	"os"
	"strconv"

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
		return
	}
	imagesWidth, imagesHeight := validateDimensions(c)

	var resizeImageInputDto dtos.ResizeImageInputDTO

	for i, image := range images {
		resizeImageInputDto.Images = append(resizeImageInputDto.Images, &dtos.ResizeImage{
			Image:  image,
			Width:  int(imagesWidth[i]),
			Height: int(imagesHeight[i]),
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

func validateDimensions(c *gin.Context) (width []int, height []int) {
	imagesWidth := c.PostFormArray("images_width")
	imagesHeight := c.PostFormArray("images_height")
	if len(imagesWidth) != len(imagesHeight) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "images_width and images_height must have the same length",
		})
		return nil, nil
	}

	parsedWidths := make([]int, len(imagesWidth))
	parsedHeights := make([]int, len(imagesHeight))

	for i, width := range imagesWidth {
		parsedWidth, err := strconv.ParseInt(width, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "images_width must be integer",
			})
			return nil, nil
		}

		parsedWidths[i] = int(parsedWidth)
		parsedHeight, err := strconv.ParseInt(imagesHeight[i], 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "images_height must be integer",
			})
			return nil, nil
		}

		parsedHeights[i] = int(parsedHeight)

		if parsedWidth <= 0 || parsedHeight <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "images_width and images_height must be greater than 0",
			})
			return nil, nil
		}
	}

	return parsedWidths, parsedHeights
}
