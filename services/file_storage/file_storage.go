package filestorage

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

type FileStorage struct{}

func (FileStorage) Store(file multipart.File, header *multipart.FileHeader) (string, error) {
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	filename := uuid.NewString() + "-" + header.Filename
	dst, err := os.Create(filepath.Join("./uploads/", filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = dst.Write(bs)
	if err != nil {
		return "", err
	}
	URL := "http://localhost:8080/pictures/" + filename + "/"
	return URL, nil
}
