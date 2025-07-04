package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mhdph/go-start/internal/app"
	"github.com/mhdph/go-start/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	app.Logger.Println("we are runing our app")

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
