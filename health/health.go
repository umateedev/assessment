package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Health struct {
	Database string `json:"database"`
	Api      string `json:"api"`
}

func HealthCheck(c echo.Context) error {
	health := Health{}
	health.Database = "ok"
	health.Api = "ok"

	return c.JSON(http.StatusOK, health)
}
