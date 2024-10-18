package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// User represents a user object.
type User struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Branch   Branch   `json:"branch"`
	Roles    []string `json:"roles"`
}

// Branch represents a branch object.
type Branch struct {
	Id int `json:"id"`
}

// Authenticate authenticates a user based on the Authorization header.
// It makes a request to the auth service to validate the token and retrieve user information.
func Authenticate(c *gin.Context) bool {
	viper.SetDefault("AUTH_API_ME", "/")

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || len(authHeader) < 7 || authHeader[0:7] != "Bearer " {
		return false
	}

	token := authHeader[7:]

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

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return false
	}

	c.Set("user", user)

	return true
}
