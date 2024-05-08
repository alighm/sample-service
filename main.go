package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	openapi "github.com/alighm/sample-service/api/gen"
	api "github.com/alighm/sample-service/api/handlers"
	"github.com/alighm/sample-service/log"
	cm "github.com/alighm/sample-service/middleware"
)

const (
	timeout = 5
)

func main() {
	log.SetLevel("info")
	logger := log.GetLogger()

	// grab env variables
	viper.AutomaticEnv()

	// api layer stitching
	UserAPIHandler := api.NewUserAPIService()
	UserAPIController := wrapMiddleware(openapi.NewUserAPIController(UserAPIHandler))

	router := openapi.NewRouter(UserAPIController)

	// creating the rest api server
	server := createServer(router)

	// channel to handle graceful shutdown of the http server (done pattern)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// starting the http server on a separate go routine
	// this way we can control the graceful shutdown of the server when a signal termination is received from main
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("failed to listen and serve: %v", err)
		}
	}()
	logger.Info("http server started: %v", time.Now())

	// the graceful shutdown of the http server
	// keep main alive until signal termination is thrown.
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// because of timeout, the context has a done pattern to close listeners
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("failed to shutdown: %v", err)
	}
	logger.Info("http server stopped: %v", time.Now())
}

func wrapMiddleware(router openapi.Router) openapi.Router {
	var routes = make(map[string]openapi.Route)

	for key, route := range router.Routes() {
		// Middlewares execute in reverse order of which they are applied,
		// So the order will be as follows
		// 1. Custom Headers
		// 2. AuthN (if applicable)
		// 3. AuthZ (if applicable)

		// every route requires Custom Headers to be handled
		route.HandlerFunc = cm.CustomHeaders(viper.GetStringSlice("APP_VERSIONS"), route.HandlerFunc)
		routes[key] = route
	}

	return openapi.NewCustomRouter(router, routes)
}

func createServer(router *mux.Router) *http.Server {
	return &http.Server{
		Addr:              ":" + strconv.Itoa(viper.GetInt("PORT")),
		ReadHeaderTimeout: timeout * time.Second,
		Handler:           router,
	}
}
