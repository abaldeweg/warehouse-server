package web

import (
	"log"
	"net/http"

	"github.com/abaldeweg/warehouse-server/products/client"
	"github.com/gin-gonic/gin"
)

type UpdateProduct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func List(c *gin.Context) {
	config, err := client.NewClient()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	results, err := config.List()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, results)
}

func Create(c *gin.Context) {
	config, err := client.NewClient()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var document client.Product
	if err := c.BindJSON(&document); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := config.Create(client.Product{
		Name:       document.Name,
		Attributes: document.Attributes,
		Variants:   document.Variants,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

    log.Println(res)

	c.Status(http.StatusCreated)
}

func Update(c *gin.Context) {
	config, err := client.NewClient()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var document UpdateProduct
	if err := c.BindJSON(&document); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := config.Update(c.Param("id"), document.Key, document.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

    log.Println(res)

	c.Status(http.StatusOK)
}

func Delete(c *gin.Context) {
	config, err := client.NewClient()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := config.Delete(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

    log.Println(res)

	c.Status(http.StatusOK)
}
