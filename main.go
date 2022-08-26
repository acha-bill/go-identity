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

	_ "go-identity/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title go-identity demo
// @version 1.0
// @description identity demo
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	router := mux.NewRouter()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// parse env vars.
	port := utils.GetEnv("PORT", "8000")
	timeout, err := time.ParseDuration(utils.GetEnv("SHUTDOWN_TIMEOUT", "5s"))
	if err != nil {
		logger.Fatal(err)
	}

	// setup routing and middlewares.
	registerHandlers(logger, router)
	registerMiddlewares(logger, router)

	// register docs
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

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
		2: {
			FirstName: "Mary",
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
