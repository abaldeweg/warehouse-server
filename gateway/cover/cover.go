package cover

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/nfnt/resize"
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
		width  uint
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

func resizeAndSaveImage(imagePath string, resizedImagePath string, width uint) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image")
	}
	defer file.Close()

	img, err := decodeImage(file, imagePath)
	if err != nil {
		return err
	}

	resizedImage := resize.Resize(width, 0, img, resize.Lanczos3)

	out, err := os.Create(resizedImagePath)
	if err != nil {
		return fmt.Errorf("failed to create resized image file")
	}
	defer out.Close()

	if err := encodeImage(out, resizedImage, imagePath); err != nil {
		return fmt.Errorf("failed to encode resized image")
	}

	return nil
}

func decodeImage(file *os.File, imagePath string) (image.Image, error) {
	ext := filepath.Ext(imagePath)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
}

func encodeImage(out *os.File, resizedImage image.Image, imagePath string) error {
	ext := filepath.Ext(imagePath)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(out, resizedImage, nil)
	case ".png":
		return png.Encode(out, resizedImage)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
}
