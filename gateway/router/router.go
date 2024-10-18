package router

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/abaldeweg/warehouse-server/gateway/proxy"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/nfnt/resize"
)

type User struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Branch   Branch   `json:"branch"`
	Roles    []string `json:"roles"`
}

type Branch struct {
	Id int `json:"id"`
}

func Routes() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Any(`/apis/core/1/*path`, func(c *gin.Context) {
		path := c.Param("path")

        if c.Request.Method == http.MethodPost {
            if match, _ := regexp.MatchString(`/api/book/cover/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`, path); match {
                if authenticate(c) {
                    re := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
                    imageUUID := re.FindString(path)

                    saveCover(c, imageUUID)

                    c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
                    return
                }
                c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
                return
            }
        }

		safePath := filepath.Join("/", path)

		if err := proxy.Proxy(c, viper.GetString("API_CORE"), safePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
			return
		}
	})

	return r
}

func saveCover(c *gin.Context, imageUUID string) {
	imageData, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload required"})
		return
	}

	imageFilename := fmt.Sprintf("%s%s", imageUUID, filepath.Ext(imageData.Filename))

	uploadsDir := "uploads"
	currentDir, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current directory"})
		return
	}

	uploadsDir = filepath.Join(currentDir, uploadsDir)

	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

	imagePath := filepath.Join(uploadsDir, imageFilename)
	if err := c.SaveUploadedFile(imageData, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	sizes := []struct {
		width  uint
		suffix string
	}{
		{400, "l"},
		{200, "m"},
		{100, "s"},
	}

	for _, size := range sizes {
		resizedImagePath := fmt.Sprintf("%s_%s%s", imageUUID, size.suffix, filepath.Ext(imageData.Filename))
		resizedImagePath = filepath.Join(uploadsDir, resizedImagePath)

		if err := resizeAndSaveImage(imagePath, resizedImagePath, size.width); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resize image"})
			return
		}
	}

	os.Remove(imagePath)
}

func resizeAndSaveImage(imagePath string, resizedImagePath string, width uint) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var img image.Image
	ext := filepath.Ext(imagePath)

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return err
	}

	resizedImage := resize.Resize(width, 0, img, resize.Lanczos3)

	out, err := os.Create(resizedImagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	switch ext {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, resizedImage, nil)
	case ".png":
		png.Encode(out, resizedImage)
	}

	return nil
}

func authenticate(c *gin.Context) bool {
	viper.SetDefault("AUTH_API_ME", "/")

	authHeader := c.GetHeader("Authorization")

	logFile, _ := os.OpenFile("/upload/auth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	logMessage := fmt.Sprintf("AUTH_API_ME: %s, authHeader: %s\n", viper.GetString("AUTH_API_ME"), authHeader)
	logFile.WriteString(logMessage)

	if authHeader == "" || len(authHeader) < 7 || authHeader[0:7] != "Bearer " {
		return false
	}

    token := authHeader[7:]

    logMessage2 := fmt.Sprintf("Token: %s\n", token)
    logFile.WriteString(logMessage2)

	req, err := http.NewRequest("GET", viper.GetString("AUTH_API_ME"), nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
