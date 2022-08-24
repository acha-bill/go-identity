package main

import (
	"context"
	"fmt"
	"go-identity/domain"
	"go-identity/handler"
	"go-identity/middleware"
	"go-identity/pkg"
	"go-identity/store"
	"go-identity/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// setup routing and middlewares.
	registerHandlers(logger, router)
	registerMiddlewares(logger, router)

	// parse env vars.
	port := utils.GetEnv("PORT", "8000")
	timeout, err := time.ParseDuration(utils.GetEnv("SHUTDOWN_TIMEOUT", "5s"))
	if err != nil {
		logger.Fatal(err)
	}

	// create and start server.
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0: %s", port),
		Handler: router,
	}
	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Printf("started on port: %s", port)

	// graceful shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = server.Shutdown(ctx)
	logger.Println("shutting down")
}

func registerHandlers(logger *log.Logger, r *mux.Router) {
	// dummy storage.
	d := map[int]*domain.Identity{
		1: {
			FirstName: "John",
			LastName:  "Doe",
			DOB:       time.Now(),
		},
	}
	identityStore := store.NewIdentity(d)
	identitySvc := pkg.NewIdentity(identityStore)
	identityHandler := handler.NewIdentity(logger, r, identitySvc)
	identityHandler.RegisterRoutes()
}

func registerMiddlewares(logger *log.Logger, r *mux.Router) {
	recovery := middleware.NewRecovery(logger)
	r.Use(recovery.Recovery)
}
