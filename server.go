package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/umateedev/assessment/database"
	"github.com/umateedev/assessment/expense"
	"github.com/umateedev/assessment/health"
)

var db *sql.DB

func main() {

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.InitDb()

	port := os.Getenv("PORT")
	log.Printf("PORT is %s", port)

	if len(port) == 0 {
		port = ":2565"
	}

	e.GET("/", landingPage)
	e.GET("/health", health.HealthCheck)

	g := e.Group("expenses")
	g.POST("", expense.CreateExpenseHandler)
	g.GET("/:id", expense.GetExpenseHandler)

	log.Printf("Server start at port %s", port)

	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal((err))
	}

	log.Printf("Server stopped")
}

func landingPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to Expenses API")
}
