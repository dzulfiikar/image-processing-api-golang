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

func TestCompressImage(t *testing.T) {
	router := gin.Default()
	routes.SetupRoutes(&router.RouterGroup)

	t.Run("Given png images with qualities request between 0-9, When compress images, Should return 200 with buffer of zip file containing compressed images", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/png-1.png", "images/png-2.png", "images/png-3.png"}
		for _, imageName := range imageNames {
			qualities, err := writer.CreateFormField("qualities")
			if err != nil {
				log.Panicln(err)
			}
	
			qualities.Write([]byte("1"))

			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/png")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		success := assert.Equal(t, 200, recorder.Code)
		if !success {
			log.Panicln(recorder.Body)
		}

		parsedZip := utils.ParseResponseToZip(recorder.Body.Bytes(), "compressed_images", "compressed_images/file.zip")


		assert.Equal(t, len(imageNames), len(parsedZip.File))
	})

	t.Run("Given png images with qualities request above 9, When compress images, Should return 400", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/png-1.png", "images/png-2.png", "images/png-3.png"}
		for _, imageName := range imageNames {
			qualities, err := writer.CreateFormField("qualities")
			if err != nil {
				log.Panicln(err)
			}
	
			qualities.Write([]byte("100"))

			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/png")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, `{"message":"qualities for png file(s) must be between 0 - 9"}`, recorder.Body.String())
	})

	t.Run("Given jpg images with qualities request between 0-100, When compress images, Should return 200 with buffer of zip file containing compressed images", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/jpg-1.jpg", "images/jpg-2.jpg", "images/jpg-3.jpg"}
		for _, imageName := range imageNames {
			qualities, err := writer.CreateFormField("qualities")
			if err != nil {
				log.Panicln(err)
			}
	
			qualities.Write([]byte("1"))

			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/jpeg")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		success := assert.Equal(t, 200, recorder.Code)
		if !success {
			log.Panicln(recorder.Body)
		}

		parsedZip := utils.ParseResponseToZip(recorder.Body.Bytes(), "compressed_images", "compressed_images/file.zip")


		assert.Equal(t, len(imageNames), len(parsedZip.File))
	})

	t.Run("Given jpg images with qualities request above 100, When compress images, Should return 400", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/jpg-1.jpg", "images/jpg-2.jpg", "images/jpg-3.jpg"}
		for _, imageName := range imageNames {
			qualities, err := writer.CreateFormField("qualities")
			if err != nil {
				log.Panicln(err)
			}
	
			qualities.Write([]byte("300"))

			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/jpeg")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, `{"message":"qualities for jpg file(s) must be between 0 - 100"}`, recorder.Body.String())
	})

	t.Run("Given images without `qualities` request, When compress images, should return 400", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/png-1.png", "images/png-2.png", "images/png-3.png"}
		for _, imageName := range imageNames {
			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/png")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, `{"message":"qualities is required"}`, recorder.Body.String())
	})

	t.Run("Given non png/jpeg images with `qualities` request, When compress images, should return 400", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(body)

		imageNames := []string{"images/webp-1.webp"}
		for _, imageName := range imageNames {
			qualities, err := writer.CreateFormField("qualities")
			if err != nil {
				log.Panicln(err)
			}
	
			qualities.Write([]byte("1"))

			imageFile, err := os.Open(imageName)
			if err != nil {
				log.Panicln(err)
			}

			defer imageFile.Close()

			imageHeader := utils.CreateFormFile("images", filepath.Base(imageFile.Name()), "image/webp")
			part, err := writer.CreatePart(imageHeader)
			if err != nil {
				log.Panicln(err)
			}

			_, err = io.Copy(part, imageFile)
			if err != nil {
				log.Panicln(err)
			}
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/compress-image", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, `{"message":"image type must be image/png or image/jpeg"}`, recorder.Body.String())
	})

}
