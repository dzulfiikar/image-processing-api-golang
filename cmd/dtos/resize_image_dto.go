package dtos

import (
	"mime/multipart"
)

type ResizeImage struct {
	Width  int
	Height int
	Image  *multipart.FileHeader
}

type ResizeImageInputDTO struct {
	Images []*ResizeImage
}

type ResizeImageOutputDTO struct {
	FileName string
	FilePath string
}
