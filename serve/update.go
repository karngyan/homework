package serve

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (s server) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	customer, err := s.ds.Get(id)
	if err != nil {
		if IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "customer not found")
		}
		return err
	}

	request := struct {
		Customer struct {
			Attributes map[string]string `json:"attributes"`
		} `json:"customer"`
	}{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	customer, err = s.ds.Update(id, request.Customer.Attributes)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct {
		Customer *Customer `json:"customer"`
	}{Customer: customer})
}
