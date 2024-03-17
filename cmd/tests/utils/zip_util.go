package utils

import (
	"archive/zip"
	"log"
	"os"
)

func ParseResponseToZip(byte []byte, targetDir string, targetPath string) *zip.ReadCloser {
	os.Mkdir(targetDir, 0777)

	_ = os.WriteFile(targetPath, byte, 0777)

	zipR, err := zip.OpenReader(targetPath)
	if err != nil {
		log.Println("Error when reading zip file")
		log.Fatal(err)
	}
	defer func() {
		zipR.Close()
		os.Remove(targetPath)
	}()

	return zipR
}
