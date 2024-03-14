package dtos

import "mime/multipart"

type CompressImage struct {
	Image   *multipart.FileHeader
	Quality int
}

type CompressImageInputDTO struct {
	CompressImages []*CompressImage
}

type CompressImageOutputDTO struct {
	FileName string
	FilePath string
}
