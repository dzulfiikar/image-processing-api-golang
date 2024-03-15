package utils

import (
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
)

func PngToJpegConverter(sourceName string, sourceImg *multipart.File, destDir string) (*os.File, error) {
	pngImg, err := png.Decode(*sourceImg)
	if err != nil {
		log.Printf("Error decoding image: %s", err)
		return nil, err
	}

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Printf("Error creating directory: %s", err)
		return nil, err
	}

	tempFile, err := os.CreateTemp(destDir, sourceName+"_*.jpg")
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return nil, err
	}

	defer tempFile.Close()

	err = jpeg.Encode(tempFile, pngImg, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Printf("Error encoding image: %s", err)
		return nil, err
	}

	return tempFile, nil
}
