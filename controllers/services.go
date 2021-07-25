package controller

import "mime/multipart"

type FileStorage interface {
	Store(file multipart.File, header *multipart.FileHeader) (URL string, filename string, err error)
	Discard(filename string) error
}
