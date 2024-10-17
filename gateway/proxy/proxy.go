package proxy

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context, serviceURL string, path string) error {
	const duration = 20 * time.Second

	ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
	defer cancel()

	req, err := request(ctx, c, serviceURL, path)
	if err != nil {
		return err
	}

	return response(req, c)
}

func request(ctx context.Context, c *gin.Context, serviceURL string, path string) (*http.Request, error) {
	fullURL := serviceURL + path + "?" + c.Request.URL.RawQuery

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, fullURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return nil, err
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}

func response(req *http.Request, c *gin.Context) error {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		} else {
			c.JSON(http.StatusBadGateway, gin.H{"msg": "Bad Gateway"})
		}
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"msg": "Bad Gateway"})
		return err
	}

	for key, values := range resp.Header {
		// Skip the Vary header if it only contains "Accept"
		if key == "Vary" && len(values) == 1 && values[0] == "Accept" {
			continue
		}

		// Remove "Accept" from the Vary header if it exists
		if key == "Vary" {
			newValues := []string{}
			for _, value := range values {
				if value != "Accept" {
					newValues = append(newValues, value)
				}
			}
			if len(newValues) > 0 {
				c.Header(key, strings.Join(newValues, ", "))
			}
			continue
		}

		c.Header(key, values[0])
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)

	return nil
}
