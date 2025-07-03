package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mhdph/go-start/internal/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	app.Logger.Println("we are runing our app")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
