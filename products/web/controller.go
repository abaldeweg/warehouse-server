package web

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/products/client"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	results, err := client.List()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
        return
	}

	c.JSON(http.StatusOK, results)
}

func Create(c *gin.Context) {
	var document client.Product
	if err := c.BindJSON(&document); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := client.Create(client.Product{
		Name:       document.Name,
		Attributes: document.Attributes,
		Variants:   document.Variants,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func Update(c *gin.Context) {
	var document UpdateProduct
	if err := c.BindJSON(&document); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
        return
	}

	res, err := client.Update(c.Param("id"), document.Key, document.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func Delete(c *gin.Context) {
	res, err := client.Delete(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
