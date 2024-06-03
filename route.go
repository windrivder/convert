package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitRoute(e *echo.Echo, cfg *Config) {
	logic := NewLogic(cfg, e.Logger)

	e.POST("/api/v1/office/convert", HandlerOfficeToJPG(logic))
}

func HandlerOfficeToJPG(logic *Logic) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req OfficeConvertRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		reply, err := logic.OfficeConvert(c.Request().Context(), &req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, reply)
	}
}
