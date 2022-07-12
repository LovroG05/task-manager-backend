package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func NewPagination(c *gin.Context) *Pagination {
	p := Pagination{}
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	sort := c.DefaultQuery("sort", "id asc")
	var err error
	p.Page, err = ParseInt(page)
	if err != nil {
		p.Page = 1
	}
	p.Limit, err = ParseInt(limit)
	if err != nil {
		p.Limit = 10
	}
	p.Sort = sort
	return &p
}

func NewDefaultPagination() *Pagination {
	p := Pagination{}
	p.Page = 1
	p.Limit = 10
	p.Sort = "id asc"
	return &p
}

func ParseInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}
