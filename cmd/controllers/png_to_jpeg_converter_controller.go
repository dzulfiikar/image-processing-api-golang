package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dzulfiikar/image-processing-api-golang/cmd/dtos"
	"github.com/dzulfiikar/image-processing-api-golang/cmd/services"
	"github.com/gin-gonic/gin"
)

/**
* Input:
* 	- Type: Multipart/form-data
* 	- Body: images: []Files
**/
func PngToJpegConverter(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		log.Println("Error parsing form : ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return
	}

	images := form.File["images"]
	if images == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "image(s) are required",
		})
		return
	}
	for i, image := range images {
		if image.Header.Get("Content-Type") != "image/png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "images type must be png",
			})
			return
		}
		images[i].Filename = "image_" + fmt.Sprintf("%d", i)
	}

	result, err := services.PngToJpegConverterService(dtos.PngToJpegInputDTO{
		Images: images,
	})

	if err != nil {
		log.Println("Error converting image : ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	defer os.Remove(result.FilePath)

	c.FileAttachment(result.FilePath, result.FileName)
}
