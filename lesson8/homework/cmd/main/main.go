package main

import (
	"homework8/internal/adapters/adrepo"
	"homework8/internal/adapters/userrepo"
	"homework8/internal/app/adapp"
	"homework8/internal/app/userapp"
	"homework8/internal/ports/httpgin"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", adapp.NewApp(adrepo.New()), userapp.NewApp(userrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
