package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/utils"
	"gocv.io/x/gocv"
)

func CompressImageService(dto dtos.CompressImageInputDTO) (dtos.CompressImageOutputDTO, error) {

	MAX_WORKERS := 10
	workers := make(chan *dtos.CompressImage, MAX_WORKERS)
	workerErrors := make(chan error)
	files := make([]*os.File, 0)
	wg := sync.WaitGroup{}

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for image := range workers {
				file, err := compressImage(image.Image, image.Quality, "compressed_images")

				if err != nil {
					workerErrors <- err
					continue
				}

				files = append(files, file)
			}
		}(i)
	}

	for _, image := range dto.CompressImages {
		workers <- image
	}
	close(workers)

	wg.Wait()

	if len(workerErrors) > 0 {
		return dtos.CompressImageOutputDTO{}, <-workerErrors
	}

	defer func() {
		for _, file := range files {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	zipFile, err := utils.ZipFiles(files, "compressed_images")
	if err != nil {
		log.Println("Error creating zip file: ", err)
		return dtos.CompressImageOutputDTO{}, err
	}

	// added sleep to handle socket timeout
	time.Sleep(100 * time.Millisecond)

	return dtos.CompressImageOutputDTO{
		FileName: "compressed_images.zip",
		FilePath: zipFile.Name(),
	}, nil

}

func compressImage(sourceImg *multipart.FileHeader, targetQuality int, targetDir string) (*os.File, error) {
	err := os.MkdirAll(targetDir, 0775)
	if err != nil {
		log.Println("Error when creating directory")
		return nil, err
	}

	file, err := sourceImg.Open()
	if err != nil {
		log.Println("Error when opening file")
		return nil, err
	}

	decodedImage, _, err := image.Decode(file)
	if err != nil {
		log.Println("Error when decoding image")
		return nil, err
	}

	sourceType := strings.Split(sourceImg.Header.Get("Content-Type"), "/")[1]
	tempFile, err := os.CreateTemp(targetDir, sourceImg.Filename+"_*."+sourceType)
	if err != nil {
		log.Println("Error when creating temp file", err)
		return nil, err
	}

	var compressionParam []int
	switch sourceType {
	case "png":
		png.Encode(tempFile, decodedImage)
		compressionParam = []int{gocv.IMWritePngCompression, targetQuality}
	case "jpeg":
		jpeg.Encode(tempFile, decodedImage, &jpeg.Options{})
		compressionParam = []int{gocv.IMWriteJpegQuality, targetQuality}
	}

	matImg := gocv.NewMat()
	matImg = gocv.IMRead(tempFile.Name(), gocv.IMReadAnyColor)
	if matImg.Empty() {
		fmt.Printf("Error reading image: %s\n", tempFile.Name())
		return nil, err
	}

	success := gocv.IMWriteWithParams(tempFile.Name(), matImg, compressionParam)
	if !success {
		log.Println("Error compressing image")
		return nil, err
	}

	return tempFile, nil
}
