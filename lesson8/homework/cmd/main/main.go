package main

import (
	"homework8/internal/adapters/adrepo"
	"homework8/internal/app"
	"homework8/internal/ports/httpfiber"
)

func main() {
	server := httpfiber.NewHTTPServer(":18080", app.NewApp(adrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
