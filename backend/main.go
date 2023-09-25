package main

import (
	"backend/controllers"
	"backend/pkg/v1/mysql"
	"backend/utils/config"
	"backend/utils/constants"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/api/health"
	"backend/router"
	"backend/scheduler"
)

func main() {
	var isDev bool
	if config.MyConfig.Environment == constants.PROD {
		isDev = false
	} else {
		isDev = true
	}

	r := router.CreateRouter(isDev)

	r.GET("/", health.HandleHealthCheck)

	r.HandleMethodNotAllowed = true
	r.NoMethod(controllers.HandleNoMethod)
	r.NoRoute(controllers.HandleNoRoutes)

	router.InitRouteV1_0_0(r)

	timeout := time.Duration(config.MyConfig.Timeout) * time.Second
	newHandler := http.TimeoutHandler(r, timeout, "Timeout!")

	server := &http.Server{
		Addr:         config.MyConfig.ServerPort,
		Handler:      newHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 25 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Start scheduler
		scheduler.StartScheduler()

		<-quit

		// close all instance
		mysql.CloseDBConnection()

		log.Println(config.MyConfig.AppName, " is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server %v: %v\n", config.MyConfig.AppName, err)
		}
		close(done)
	}()

	log.Println(config.MyConfig.AppName, " is ready to handle requests at", config.MyConfig.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", config.MyConfig.ServerPort, err)
	}

	<-done
	log.Println(config.MyConfig.AppName, " stopped")
}
