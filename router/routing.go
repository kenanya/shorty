package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	v1 "github.com/kenanya/shorty/controllers/v1"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func AssignRouting(e *echo.Echo, db *mongo.Database) *echo.Echo {

	//api v1 handlers
	api := v1.Controller{
		CurDB: db,
	}

	g := e.Group("/v1")

	g.GET("/index", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})
	g.POST("/index2", func(c echo.Context) error {
		time.Sleep(50 * time.Second)
		return c.JSON(http.StatusOK, true)
	})

	g.POST("/shorten", api.CreateShortenURL)
	g.GET("/:shortcode", api.GetURLByShortCode)
	g.GET("/:shortcode/stats", api.GetURLStatByShortCode)

	fmt.Printf("## TGFWM, port : <%+v>\n\n", e.Server.Addr)

	// Start server
	go func() {
		if err := e.Start(":" + e.Server.Addr); err != nil {
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
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return e
}
