package dtos

import (
	"mime/multipart"
)

type ImageDimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

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
