package services

import (
	"log"
	"mime/multipart"
	"os"
	"sync"

	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/dtos"
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/utils"
)

func PngToJpegConverterService(dto dtos.PngToJpegInputDTO) (dtos.PngToJpegOutputDTO, error) {
	MAX_WORKERS := 10
	workers := make(chan *multipart.FileHeader, MAX_WORKERS)
	workerErrors := make(chan error)
	files := make([]*os.File, 0)
	wg := sync.WaitGroup{}

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for image := range workers {
				imgReader, err := image.Open()
				if err != nil {
					workerErrors <- err
					continue
				}
				file, err := utils.PngToJpegConverter(image.Filename, &imgReader, "converted_images")

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

	zipFile, err := utils.ZipFiles(files, "converted_images")
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
