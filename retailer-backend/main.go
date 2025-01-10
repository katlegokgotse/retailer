package main

import (
	"log"
	"net/http"
	"main.go/router"
	"main.go/models"
)

func main() {
	http.HandleFunc("/users") //HandleFunc registers the handler function for the given pattern in [DefaultServeMux]. The documentation for [ServeMux] explains how patterns are matched.
	var users = []models.User{
		{ID: "1", Username: "katlegokgotse", Email: "katlegokgotse88@gmail.com", Password: "123456", PhoneNumber: "0662362301"},
	}
	router.authRouter(controllers)
	log.Fatal(http.ListenAndServe(":8080", nil)) //ListenAndServe listens on the TCP network address addr and then calls [Serve] with handler to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives
}
