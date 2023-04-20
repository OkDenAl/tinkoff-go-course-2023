package main

import (
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app/adsapp"
	"homework9/internal/app/userapp"
	"homework9/internal/ports/httpgin"
)

func main() {
	userRepo := userrepo.New()
	server := httpgin.NewHTTPServer(":18080", adsapp.NewApp(adrepo.New(), userRepo), userapp.NewApp(userRepo))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
