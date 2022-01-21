package serve

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (s server) Get(c echo.Context) error {

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

	response := struct {
		Customer *Customer `json:"customer"`
	}{customer}

	return c.JSON(http.StatusOK, response)
}
