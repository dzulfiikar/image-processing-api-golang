package dtos

import (
	"mime/multipart"
)

type PngToJpegInputDTO struct {
	Images []*multipart.FileHeader
}

type PngToJpegOutputDTO struct {
	FileName string
	FilePath string
}
