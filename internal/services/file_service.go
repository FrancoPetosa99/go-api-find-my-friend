package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-api-find-my-friend/pkg/errors"
)

type FileService struct {
	uploadPath string
	maxSize    int64
}

func NewFileService() *FileService {
	return &FileService{
		uploadPath: "./uploads/pets",
		maxSize:    10 * 1024 * 1024, // 10MB
	}
}

func (s *FileService) UploadPetImage(file *multipart.FileHeader) (string, error) {
	// Validar tipo de archivo
	if !s.isValidImageType(file.Header.Get("Content-Type")) {
		return "", errors.NewBadRequestError("Invalid file type. Only images are allowed")
	}

	// Validar tamaño
	if file.Size > s.maxSize {
		return "", errors.NewBadRequestError("File too large. Maximum size is 10MB")
	}

	// Crear directorio si no existe
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return "", errors.NewInternalServerError("Failed to create upload directory")
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("pet_%s%s", timestamp, ext)
	filepath := filepath.Join(s.uploadPath, filename)

	// Guardar archivo
	if err := s.saveFile(file, filepath); err != nil {
		return "", errors.NewInternalServerError("Failed to save file")
	}

	// Retornar URL relativa
	return fmt.Sprintf("/uploads/pets/%s", filename), nil
}

func (s *FileService) isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if strings.Contains(contentType, validType) {
			return true
		}
	}
	return false
}

func (s *FileService) saveFile(file *multipart.FileHeader, filepath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copiar contenido del archivo
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}
		if _, err := dst.Write(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileService) DeleteFile(filepath string) error {
	if filepath == "" {
		return nil
	}

	// Convertir URL relativa a path absoluto
	if strings.HasPrefix(filepath, "/uploads/") {
		filepath = "." + filepath
	}

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
