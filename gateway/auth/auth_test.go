package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	mockAuthService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := User{
			Id:       1,
			Username: "testuser",
			Branch: Branch{
				Id: 1,
			},
			Roles: []string{"admin"},
		}
		json.NewEncoder(w).Encode(user)
	}))
	defer mockAuthService.Close()

	viper.Set("AUTH_API_ME", mockAuthService.URL)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/", nil)

	testCases := []struct {
		name          string
		authorization string
		expected      bool
	}{
		{"Missing Authorization header", "", false},
		{"Invalid Authorization header", "InvalidToken", false},
		{"Valid Authorization header", "Bearer test-token", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c.Request.Header.Set("Authorization", tc.authorization)
			assert.Equal(t, tc.expected, Authenticate(c))
		})
	}
}

func TestUser(t *testing.T) {
    user := User{
        Id:       1,
        Username: "testuser",
        Branch: Branch{
            Id: 1,
        },
        Roles: []string{"admin"},
    }

    mockAuthService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") != "Bearer test-token" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        json.NewEncoder(w).Encode(user)
    }))
    defer mockAuthService.Close()

    viper.Set("AUTH_API_ME", mockAuthService.URL)

    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Request, _ = http.NewRequest("GET", "/", nil)
    c.Request.Header.Set("Authorization", "Bearer test-token")

    Authenticate(c)

    assert.Equal(t, user, c.MustGet("user").(User))
}
