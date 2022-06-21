package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/simonsanchez/echo-context/token"
)

func withCustomContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{
			Context: c,
			Store:   make(map[string]string),
		}
		return next(cc)
	}
}

func withDefaultTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := token.Default()
		if err != nil {
			return err
		}

		cc := c.(*CustomContext)
		cc.Token = token
		return next(cc)
	}
}

func withCurrentTimeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*CustomContext)

		cc.Store["key"] = "value"
		cc.Store["foo"] = "bar"
		cc.Store["now"] = time.Now().String()

		return next(cc)
	}
}
