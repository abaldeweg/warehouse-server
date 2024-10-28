package cover

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSaveCover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		SaveCover(c, "36ee6d5c-820b-4f0c-9637-73b63dacc2a7")
	})

	testCases := []struct {
		imageName      string
		expectedImages []string
	}{
		{"test.jpg", []string{
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-l.jpg",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-m.jpg",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-s.jpg",
		}},
		{"test.png", []string{
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-l.png",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-m.png",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-s.png",
		}},
		{"test.webp", []string{
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-l.webp",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-m.webp",
			"36ee6d5c-820b-4f0c-9637-73b63dacc2a7-s.webp",
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.imageName, func(t *testing.T) {
			imageData, err := os.ReadFile(tc.imageName)
			if err != nil {
				t.Fatalf("Error reading image file: %v", err)
			}

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("cover", tc.imageName)
			_, _ = io.Copy(part, bytes.NewReader(imageData))
			writer.Close()

			req, _ := http.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			currentDir, _ := os.Getwd()

			for _, expectedImage := range tc.expectedImages {
				expectedFilePath := filepath.Join(currentDir, uploadsDir, expectedImage)
				if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s to exist", expectedFilePath)
				} else {
					os.Remove(expectedFilePath)
				}
			}

			assert.NoFileExists(t, filepath.Join(currentDir, uploadsDir, "36ee6d5c-820b-4f0c-9637-73b63dacc2a7.jpg"))
			assert.NoFileExists(t, filepath.Join(currentDir, uploadsDir, "36ee6d5c-820b-4f0c-9637-73b63dacc2a7.png"))
			assert.NoFileExists(t, filepath.Join(currentDir, uploadsDir, "36ee6d5c-820b-4f0c-9637-73b63dacc2a7.webp"))

			os.RemoveAll(filepath.Join(currentDir, uploadsDir))
		})
	}
}
