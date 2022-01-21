package serve

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func isTimestamp(value string) bool {
	num, err := strconv.Atoi(value)

	if err != nil {
		return false
	}

	return num > 99999999 && num < 10000000000
}

func (s server) Create(c echo.Context) error {
	request := struct {
		Customer struct {
			ID         int               `json:"id"`
			Attributes map[string]string `json:"attributes"`
		} `json:"customer"`
	}{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	if val, ok := request.Customer.Attributes["email"]; !ok || val == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email attribute is required")
	}

	if !isTimestamp(request.Customer.Attributes["created_at"]) {
		request.Customer.Attributes["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	}

	customer, err := s.ds.Create(request.Customer.ID, request.Customer.Attributes)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		Customer *Customer `json:"customer"`
	}{Customer: customer})
}
