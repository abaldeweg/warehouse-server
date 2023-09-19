package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/incwadi-warehouse/monorepo-go/security/authentication"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Me(c *gin.Context) {
	if c.GetHeader("Authorization") == "Bearer token" {
		u := authentication.User{}
		u.Id = 1
		u.Username = "admin"
		u.Branch.Id = 1
		u.Roles = []string{"ROLE_ADMIN", "ROLE_USER"}

		c.JSON(200, u)
		return
	}

	c.JSON(401, Response{401, "UNAUTHORIZED"})
}
