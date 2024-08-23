package main

import (
	"net/http"

	"friendsocial/postgres"
	"friendsocial/users"
)

func main() {
	postgres.InitDB()
	defer postgres.CloseDB()

	mux := http.NewServeMux()

	userServices := users.NewService(postgres.DB)
	userManager := users.NewUserHTTPHandler(userServices)

	mux.HandleFunc("POST /users", userManager.HandleHTTPPost)
	mux.HandleFunc("GET /users", userManager.HandleHTTPGet)
	mux.HandleFunc("GET /users/{id}", userManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /users/{id}", userManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /users/{id}", userManager.HandleHTTPDelete)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
