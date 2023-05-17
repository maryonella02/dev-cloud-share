package main

import (
	"auth/controllers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"auth/services"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	authService := services.NewAuthService(client)
	authController := controllers.NewAuthController(authService)
	router := mux.NewRouter()
	authController.RegisterRoutes(router)

	// Set up server
	server := &http.Server{
		Addr:    ":8083",
		Handler: router,
	}

	// Start server in a separate goroutine
	go func() {
		fmt.Println("Starting server on port 8083")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shut down the server gracefully with a timeout
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server gracefully stopped")
}
