package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/utils"

	"gocv.io/x/gocv"
)

func ResizeImageService(dto dtos.ResizeImageInputDTO) (dtos.ResizeImageOutputDTO, error) {
	MAX_WORKERS := 10
	workers := make(chan *dtos.ResizeImage, MAX_WORKERS)
	workerErrors := make(chan error)
	files := make([]*os.File, 0)
	wg := sync.WaitGroup{}

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for image := range workers {
				file, err := resizeImage(image)

				if err != nil {
					workerErrors <- err
					continue
				}

				files = append(files, file)
			}
		}(i)
	}

	for _, image := range dto.Images {
		workers <- image
	}
	close(workers)

	wg.Wait()

	if len(workerErrors) > 0 {
		return dtos.ResizeImageOutputDTO{}, <-workerErrors
	}

	defer func() {
		for _, file := range files {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	zipFile, err := utils.ZipFiles(files, "resized_images")
	if err != nil {
		log.Printf("Error creating zip file: %s", err)
		return dtos.ResizeImageOutputDTO{}, err
	}

	return dtos.ResizeImageOutputDTO{
		FileName: "resized_images.zip",
		FilePath: zipFile.Name(),
	}, nil
}

func resizeImage(sourceImg *dtos.ResizeImage) (*os.File, error) {

	imgFile, err := sourceImg.Image.Open()
	if err != nil {
		fmt.Printf("Error opening image: %s", err)
		return nil, err
	}

	decodedImage, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Error when decoding image")
		return nil, err
	}

	err = os.Mkdir("resized_images", 0777)
	if err != nil {
		log.Println("Error when creating dir")
		return nil, err
	}

	sourceType := strings.Split(sourceImg.Image.Header.Get("Content-Type"), "/")[1]
	tempFile, err := os.CreateTemp("resized_images", sourceImg.Image.Filename+"_*."+sourceType)
	if err != nil {
		log.Println("Error when creating temp file", err)
		return nil, err
	}

	switch sourceType {
	case "png":
		png.Encode(tempFile, decodedImage)
	case "jpeg":
		jpeg.Encode(tempFile, decodedImage, &jpeg.Options{Quality: 100})
	}

	matImg := gocv.NewMat()
	defer matImg.Close()

	matImg = gocv.IMRead(tempFile.Name(), gocv.IMReadAnyColor)
	if matImg.Empty() {
		fmt.Printf("Error reading image: %s\n", tempFile.Name())
		return nil, err
	}

	gocv.Resize(matImg, &matImg, image.Point{
		sourceImg.Width,
		sourceImg.Height,
	}, float64(0), float64(0), gocv.InterpolationDefault)

	gocv.IMWrite(tempFile.Name(), matImg)

	return tempFile, nil
}
