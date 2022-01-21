package serve

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (s server) Delete(c echo.Context) error {

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

	if err := s.ds.Delete(customer.ID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
