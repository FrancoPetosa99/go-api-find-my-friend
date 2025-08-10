package storage_provider

import (
	"context"
	"fmt"
	"go-api-find-my-friend/pkg/config"
	"log"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

var (
	cloudinaryClientInstance *CloudinaryClient
	cloudinaryClientOnce     sync.Once
)

func NewCloudinary() *CloudinaryClient {
	cloudinaryClientOnce.Do(func() {
		cld, err := cloudinary.NewFromParams(config.ConfigInstance.Cloudinary.CloudName, config.ConfigInstance.Cloudinary.APIKey, config.ConfigInstance.Cloudinary.APISecret)
		if err != nil {
			log.Fatalf("Error creating Cloudinary client: %v", err)
		}
		cloudinaryClientInstance = &CloudinaryClient{
			cld: cld,
		}
	})
	return cloudinaryClientInstance
}

func (c *CloudinaryClient) Upload(file *multipart.FileHeader) (string, error) {
	uploadResult, err := c.cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		ResourceType: "auto",
		Folder:       "find-my-friend",
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (c *CloudinaryClient) Delete(fileURL string) error {
	publicID, err := extractPublicIDFromURL(fileURL)
	if err != nil {
		log.Printf("Failed to extract public ID from URL %s: %v", fileURL, err)
		return err
	}

	result, err := c.cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
	fmt.Println(result)
	if err != nil {
		log.Printf("Failed to delete file with public ID %s: %v", publicID, err)
		return err
	}

	log.Printf("Successfully deleted file with public ID: %s", publicID)
	return nil
}

func extractPublicIDFromURL(url string) (string, error) {
	// Según la documentación de Cloudinary:
	// URL format: https://res.cloudinary.com/cloud_name/image/upload/v1234567890/folder/filename.jpg
	// Public ID: folder/filename (sin extensión)

	uploadIndex := strings.Index(url, "/upload/")
	if uploadIndex == -1 {
		return "", fmt.Errorf("invalid Cloudinary URL format")
	}

	// Obtener la parte después de /upload/
	pathAfterUpload := url[uploadIndex+8:] // 8 es la longitud de "/upload/"

	// Dividir por "/" para separar las partes
	parts := strings.Split(pathAfterUpload, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid Cloudinary URL format")
	}

	// El public_id es la ruta completa sin la extensión
	// Ejemplo: "find-my-friend/pets/1234567890.jpg" -> "find-my-friend/pets/1234567890"

	// Construir el public_id (todos los segmentos excepto el último si es una extensión)
	var publicIDParts []string

	// Saltar el primer segmento si es un número de versión (v1234567890)
	startIndex := 0
	if len(parts) > 0 && len(parts[0]) > 1 && parts[0][0] == 'v' && isNumeric(parts[0][1:]) {
		startIndex = 1
	}

	// Agregar todos los segmentos excepto la extensión del último
	for i := startIndex; i < len(parts); i++ {
		if i == len(parts)-1 {
			// Para el último segmento, remover la extensión
			lastPart := parts[i]
			dotIndex := strings.LastIndex(lastPart, ".")
			if dotIndex != -1 {
				lastPart = lastPart[:dotIndex]
			}
			publicIDParts = append(publicIDParts, lastPart)
		} else {
			publicIDParts = append(publicIDParts, parts[i])
		}
	}

	if len(publicIDParts) == 0 {
		return "", fmt.Errorf("invalid Cloudinary URL format")
	}

	publicID := strings.Join(publicIDParts, "/")
	return publicID, nil
}

func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
