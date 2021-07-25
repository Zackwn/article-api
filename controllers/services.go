package controller

import "mime/multipart"

type FileStorage interface {
	Store(file multipart.File, header *multipart.FileHeader) (string, error)
}
