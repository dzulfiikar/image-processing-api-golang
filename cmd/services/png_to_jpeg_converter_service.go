package services

import (
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"archive/zip"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
)

var MAX_WORKERS = 10

func PngToJpegConverterService(dto dtos.PngToJpegInputDTO) (dtos.PngToJpegOutputDTO, error) {

	workers := make(chan *multipart.FileHeader, MAX_WORKERS)
	workerErrors := make(chan error)
	files := make([]*os.File, 0)
	wg := sync.WaitGroup{}

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for image := range workers {
				file, err := pngToJpegConverter(image)

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
		return dtos.PngToJpegOutputDTO{}, <-workerErrors
	}

	zipFile, err := zipFiles(files)
	if err != nil {
		log.Printf("Error creating zip file: %s", err)
		return dtos.PngToJpegOutputDTO{}, err
	}

	defer func() {
		for _, file := range files {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	return dtos.PngToJpegOutputDTO{
		FileName: "converted_images.zip",
		FilePath: zipFile.Name(),
	}, nil
}

func pngToJpegConverter(image *multipart.FileHeader) (*os.File, error) {
	reader, err := image.Open()
	if err != nil {
		log.Printf("Error opening image: %s", err)
		return nil, err
	}

	defer reader.Close()
	pngImg, err := png.Decode(reader)
	if err != nil {
		log.Printf("Error decoding image: %s", err)
		return nil, err
	}

	// io.Writer
	file, err := os.CreateTemp("converted_images", image.Filename+"_*.jpg")
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return nil, err
	}

	defer func() {
		file.Close()
	}()

	err = jpeg.Encode(file, pngImg, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Printf("Error encoding image: %s", err)
		return nil, err
	}

	return file, nil
}

func zipFiles(files []*os.File) (*os.File, error) {
	tempFileName := "converted_images" + time.DateOnly
	f, err := os.CreateTemp("converted_images", tempFileName)
	if err != nil {
		log.Printf("Error creating file: %s", err)
	}

	defer func() {
		f.Close()
	}()

	writer := zip.NewWriter(f)
	defer writer.Close()

	for _, file := range files {
		f, err := writer.Create(file.Name())
		if err != nil {
			log.Printf("Error creating file: %s", err)
		}
		fileBytes, _ := os.ReadFile(file.Name())
		_, err = f.Write(fileBytes)
		if err != nil {
			log.Printf("Error writing file: %s", err)
		}
	}

	return f, nil
}
