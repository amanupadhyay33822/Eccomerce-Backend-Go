package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationParams struct {
	Page   int
	Limit  int
	Skip   int
	Filter bson.M
}

func GetPaginationParams(c *gin.Context) *PaginationParams {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	return &PaginationParams{
		Page:  page,
		Limit: limit,
		Skip:  skip,
	}
}

func GetProductFilter(c *gin.Context) bson.M {
	filter := bson.M{}

	category := c.Query("category")
	if category != "" {
		filter["category"] = category
	}

	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	if minPrice != "" || maxPrice != "" {
		priceFilter := bson.M{}
		if minPrice != "" {
			if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
				priceFilter["$gte"] = min
			}
		}
		if maxPrice != "" {
			if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
				priceFilter["$lte"] = max
			}
		}
		if len(priceFilter) > 0 {
			filter["price"] = priceFilter
		}
	}

	return filter
}

func GetFindOptions(pagination *PaginationParams) *options.FindOptions {
	return options.Find().SetSkip(int64(pagination.Skip)).SetLimit(int64(pagination.Limit))
}
