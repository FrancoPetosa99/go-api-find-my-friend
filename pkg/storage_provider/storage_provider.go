package storage_provider

import "mime/multipart"

type StorageProvider interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(fileURL string) error
}
