package services

import (
	"fmt"
	"image"
	"log"
	"os"
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

	zipFile, err := utils.ZipFiles(files, "compressed_resized_images")
	if err != nil {
		log.Printf("Error creating zip file: %s", err)
		return dtos.ResizeImageOutputDTO{}, err
	}

	return dtos.ResizeImageOutputDTO{
		FileName: "compressed_resized_images.zip",
		FilePath: zipFile.Name(),
	}, nil
}

func resizeImage(sourceImg *dtos.ResizeImage) (*os.File, error) {

	imgReader, err := sourceImg.Image.Open()
	if err != nil {
		fmt.Printf("Error opening image: %s", err)
		return nil, err
	}

	jpegImg, err := utils.PngToJpegConverter(sourceImg.Image.Filename, &imgReader, "compressed_resized_images")
	if err != nil {
		fmt.Printf("Error converting image: %s", err)
		return nil, err
	}

	img := gocv.NewMat()
	defer img.Close()

	img = gocv.IMRead(jpegImg.Name(), gocv.IMReadAnyColor)
	if img.Empty() {
		fmt.Printf("Error reading image: %s\n", jpegImg.Name())
		return nil, err
	}

	gocv.Resize(img, &img, image.Point{sourceImg.Width, sourceImg.Height}, float64(0), float64(0), gocv.IMWriteJpegQuality)

	// defer func() {
	// 	jpegImg.Close()
	// 	os.Remove(jpegImg.Name())
	// }()

	resizedImg := gocv.IMWrite(jpegImg.Name(), img)
	if !resizedImg {
		fmt.Println("Error writing image")
		return nil, err
	}

	return jpegImg, nil
}
