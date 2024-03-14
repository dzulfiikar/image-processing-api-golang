package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/services"
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

	images := validateImages(form, c)
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

func validateImages(form *multipart.Form, c *gin.Context) []*multipart.FileHeader {
	images := form.File["images"]
	if images == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "image(s) are required",
		})
		return nil
	}

	for i, image := range images {
		if image.Header.Get("Content-Type") != "image/png" && image.Header.Get("Content-Type") != "image/jpeg" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "image type must be image/png or image/jpeg",
			})
			return nil
		}

		images[i].Filename = "image_" + fmt.Sprintf("%d", i)
	}

	return images
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
