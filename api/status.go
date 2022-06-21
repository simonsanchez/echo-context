package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Status(c echo.Context) error {
	cc := c.(*CustomContext)

	scopes, err := cc.Token.Scopes()
	if err != nil {
		return err
	}

	response := map[string]any{
		"status": "up",
		"store":  cc.Store,
		"scopes": scopes,
	}

	return c.JSON(http.StatusOK, response)
}
