package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/abaldeweg/warehouse-server/gateway/db/mdb"
	"github.com/gin-gonic/gin"
)

// AnalyzeController struct for analyze controller.
type AnalyzeController struct {
	repo *repository.AnalyzeRepository
}

// NewAnalyzeController creates a new analyze controller.
func NewAnalyzeController(db *mdb.MDBClient) *AnalyzeController {
	return &AnalyzeController{
		repo: repository.NewAnalyzeRepository(db),
	}
}

// Create handles the creation of analyze data.
func (ac *AnalyzeController) Create(c *gin.Context) {
	analyzeShopSearch := models.AnalyzeShopSearch{
		Term:   "",
		Branch: 0,
		Genre:  0,
		Page:   1,
		Date:   time.Now().Format("2006-01-02 15:04:05"),
	}

	options := c.Query("options")
	var opts models.AnalyzeShopSearchOptions
	if options != "" {
		if err := json.Unmarshal([]byte(options), &opts); err != nil {
			c.Next()
			return
		}
	}

	if opts.Term != "" {
		analyzeShopSearch.Term = opts.Term
	}

	getFilterValue := func(filters []models.AnalyzeShopSearchFilter, field string) int {
		for _, f := range filters {
			if f.Field == field {
				switch v := f.Value.(type) {
				case string:
					n, _ := strconv.Atoi(v)
					return n
				}
			}
		}
		return 0
	}

	if b := getFilterValue(opts.Filter, "branch"); b != 0 {
		analyzeShopSearch.Branch = b
	}

	if g := getFilterValue(opts.Filter, "genre"); g != 0 {
		analyzeShopSearch.Genre = g
	}

	if opts.Offset != 0 {
		analyzeShopSearch.Page = opts.Offset/20 + 1
	}

	ac.repo.Add(analyzeShopSearch)
	c.Next()
}
