package tests

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/routes"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPngToJpegConverter(t *testing.T) {
	router := gin.Default()
	routes.SetupRoutes(&router.RouterGroup)

	t.Run("Given png images, When convert png to jpeg, Should return 200 with buffer of zip file containing converted images", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		pngImages := []string{"images/png-1.png", "images/png-2.png", "images/png-3.png"}

		for _, pngImage := range pngImages {
			imageFile, err := os.Open(pngImage)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(pngImage), "image/png")
			part, err := writer.CreatePart(imageHeader)

			if err != nil {
				log.Panicln(err)
			}
			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}

		}

		err := writer.Close()
		if err != nil {
			log.Panicln(err)
		}

		req, _ := http.NewRequest("POST", "/png-to-jpeg", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)

		parsedZip := utils.ParseResponseToZip(recorder.Body.Bytes(), "converted_images","file.zip")

		assert.Equal(t, len(pngImages), len(parsedZip.File))
	})

	t.Run("Given non png images, When convert png to jpeg, Should return 400", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		jpgImage := "images/jpg-2.jpg"
		imageFile, err := os.Open(jpgImage)
		if err != nil {
			panic("Error when reading file")
		}

		defer imageFile.Close()

		imageHeader := utils.CreateFormFile("images", filepath.Base(jpgImage), "image/jpeg")
		part, err := writer.CreatePart(imageHeader)

		if err != nil {
			log.Panicln(err)
		}
		_, err = io.Copy(part, imageFile)
		if err != nil {
			log.Panicln(err)
		}

		err = writer.Close()
		if err != nil {
			log.Panicln(err)
		}

		req, _ := http.NewRequest("POST", "/png-to-jpeg", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
}
