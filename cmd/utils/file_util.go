package utils

import (
	"archive/zip"
	"log"
	"os"
	"time"
)

func ZipFiles(files []*os.File, targetDir string) (*os.File, error) {
	tempFileName := targetDir + time.DateOnly
	f, err := os.CreateTemp(targetDir, tempFileName)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return nil, err
	}

	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	for _, file := range files {
		f, err := writer.Create(file.Name())
		if err != nil {
			log.Printf("Error creating file: %s", err)
			return nil, err
		}
		fileBytes, _ := os.ReadFile(file.Name())
		_, err = f.Write(fileBytes)
		if err != nil {
			log.Printf("Error writing file: %s", err)
			return nil, err
		}
	}

	return f, nil
}
