package cover

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"golang.org/x/image/webp"
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
		resizedImagePath := filepath.Join(uploadsDir, fmt.Sprintf("%s-%s%s", imageUUID, size.suffix, ".jpg"))

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
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	var img image.Image

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	mimeType := http.DetectContentType(fileHeader)
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	switch mimeType {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	case "image/webp":
		img, err = convertWebp(file)
	default:
		return fmt.Errorf("unsupported image format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	originalBounds := img.Bounds()
	aspectRatio := float64(originalBounds.Dx()) / float64(originalBounds.Dy())
	height := int(float64(width) / aspectRatio)

	resizedImage := imaging.Resize(img, width, height, imaging.Lanczos)

	outFile, err := os.Create(resizedImagePath)
	if err != nil {
		return fmt.Errorf("failed to create resized image file: %w", err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImage, nil)
	if err != nil {
		return fmt.Errorf("failed to encode resized image: %w", err)
	}

	return nil
}

func convertWebp(file *os.File) (image.Image, error) {
  var err error
  var img image.Image
	img, err = webp.Decode(file)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	img, err = jpeg.Decode(buf)
	if err != nil {
		return nil, err
	}

	return img, nil
}
