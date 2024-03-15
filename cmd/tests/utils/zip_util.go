package utils

import (
	"archive/zip"
	"log"
	"os"
)

func ParseResponseToZip(byte []byte, fullPath string) *zip.ReadCloser {
	_ = os.WriteFile(fullPath, byte, 0777)

	zipR, err := zip.OpenReader(fullPath)
	if err != nil {
		log.Println("Error when reading zip file")
		log.Fatal(err)
	}
	defer func() {
		zipR.Close()
		os.Remove(fullPath)
	}()

	return zipR
}
