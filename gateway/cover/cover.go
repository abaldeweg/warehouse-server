package cover

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

const uploadsDir = "uploads"

// SaveCover saves the uploaded cover image in different sizes.
func SaveCover(c *gin.Context, imageUUID string) {
	imageData, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload required"})
		return
	}

	if err := saveResizedImages(c, imageData, imageUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusOK)
}

func saveResizedImages(c *gin.Context, imageData *multipart.FileHeader, imageUUID string) error {
	imagePath, err := saveUploadedImage(c, imageData, imageUUID)
	if err != nil {
		return fmt.Errorf("failed to save uploaded image: %w", err)
	}
	defer os.Remove(imagePath)

	sizes := []struct {
		width  int
		suffix string
	}{
		{400, "l"},
		{200, "m"},
		{100, "s"},
	}

	for _, size := range sizes {
		resizedImagePath := filepath.Join(uploadsDir, fmt.Sprintf("%s-%s%s", imageUUID, size.suffix, filepath.Ext(imageData.Filename)))

		if err := resizeAndSaveImage(imagePath, resizedImagePath, size.width); err != nil {
			return fmt.Errorf("failed to resize image: %w", err)
		}
	}

	return nil
}

func saveUploadedImage(c *gin.Context, imageData *multipart.FileHeader, imageUUID string) (string, error) {
	imageFilename := fmt.Sprintf("%s%s", imageUUID, filepath.Ext(imageData.Filename))
	currentDir, _ := os.Getwd()
	uploadsDirPath := filepath.Join(currentDir, uploadsDir)

	if err := os.MkdirAll(uploadsDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory")
	}

	imagePath := filepath.Join(uploadsDirPath, imageFilename)
	if err := c.SaveUploadedFile(imageData, imagePath); err != nil {
		return "", fmt.Errorf("failed to save image")
	}

	return imagePath, nil
}

func resizeAndSaveImage(imagePath string, resizedImagePath string, width int) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	resizedImage := imaging.Resize(img, width, 0, imaging.Lanczos)

	if err := imaging.Save(resizedImage, resizedImagePath); err != nil {
		return fmt.Errorf("failed to save resized image: %w", err)
	}

	return nil
}
