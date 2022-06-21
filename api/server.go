package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/simonsanchez/echo-context/token"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	// add services here
}

type CustomContext struct {
	echo.Context

	Store map[string]string

	Token token.Token
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

	return <-shutdown
}

func applyMiddleware(e *echo.Echo) {
	// This middleware must come first!
	e.Use(withCustomContextMiddleware)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add JWT to context
	e.Use(withDefaultTokenMiddleware)

	// Sample middleware using CustomContext
	e.Use(withCurrentTimeMiddleware)
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
		"error": message,
	})
}
