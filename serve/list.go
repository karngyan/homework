package serve

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (s server) List(c echo.Context) error {

	page := 1
	perPage := 25

	if val, err := strconv.Atoi(c.QueryParam("page")); err == nil && val > 0 {
		page = val
	}
	if val, err := strconv.Atoi(c.QueryParam("per_page")); err == nil && val > 0 {
		perPage = val
	}

	total, err := s.ds.TotalCustomers()
	if err != nil {
		return err
	}

	customers, err := s.ds.List(page, perPage)
	if err != nil {
		return err
	}

	reply := struct {
		Customers []*Customer `json:"customers"`
		Meta      struct {
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
			Total   int `json:"total"`
		} `json:"meta"`
	}{}

	reply.Meta.Total = total

	if len(customers) == 0 {
		reply.Meta.Page = 1
	} else {
		reply.Meta.Page = page
	}
	reply.Meta.PerPage = perPage

	reply.Customers = customers

	return c.JSON(http.StatusOK, reply)
}
