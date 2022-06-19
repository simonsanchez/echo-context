package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	// add services here
}

func (s *Server) Listen(port string) error {
	e := echo.New()

	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.SetLevel(log.INFO)

	applyMiddleware(e)
	applyRoutes(e, s)

	shutdown := make(chan error)
	go gracefulShutdown(e, shutdown)

	err := e.Start(":" + port)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	return nil
}

func applyMiddleware(e *echo.Echo) {
	// add echo context extension FIRST!

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}

func applyRoutes(e *echo.Echo, s *Server) {
	e.GET("/status", s.Status)
}

func gracefulShutdown(e *echo.Echo, c chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	e.Logger.Info("shutting down server: " + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := e.Shutdown(ctx)
	if err != nil {
		c <- err
	} else {
		c <- nil
	}
}

func (s *Server) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "up"})
}

/*
Endpoints can return a &CustomError directly and we can branch
here to populate the different error messages, e.g. all the db errors.
*/

const (
	defaultErrorMessage   = "unexpected error occurred"
	defaultHttpStatusCode = http.StatusInternalServerError
)

func customHTTPErrorHandler(err error, c echo.Context) {
	message, code := defaultErrorMessage, defaultHttpStatusCode

	switch err.(type) {
	case *echo.HTTPError:
		// we don't override `message` to avoid leaking sensitive errors
		c.Logger().Error(err)
		code = err.(*echo.HTTPError).Code
	case *CustomError:
		e := err.(*CustomError)
		c.Logger().Error(e.Internal)
		code = e.Code
		message = e.Public
	default:
		c.Logger().Error(err)

		// Additional db error checks here
	}

	c.JSON(code, map[string]string{
		"errors": message,
	})
}
