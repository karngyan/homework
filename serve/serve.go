package serve

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type server struct {
	ds Datastore
}

func ListenAndServe(address string, datastore Datastore) error {

	log.Println("listening on", address)

	s := server{ds: datastore}

	// Setup
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/customers", s.List)
	e.POST("/customers", s.Create)
	e.GET("/customers/:id", s.Get)
	e.PATCH("/customers/:id", s.Update)
	e.DELETE("/customers/:id", s.Delete)

	// Start server
	go func() {
		if err := e.Start(address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}
