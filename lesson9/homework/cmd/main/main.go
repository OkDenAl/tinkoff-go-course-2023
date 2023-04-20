package main

import "homework9/internal/ports/httpgin"

func main() {
	userRepo := userrepo.New()
	server := httpgin.NewHTTPServer(":18080", adapp.NewApp(adrepo.New(), userRepo), userapp.NewApp(userRepo))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
