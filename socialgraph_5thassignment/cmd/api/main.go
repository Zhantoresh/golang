package main

import (
	"log"
	"net/http"

	"socialgraph_5thassignment/internal/app"
	"socialgraph_5thassignment/internal/handler"
	"socialgraph_5thassignment/internal/repository"
)

func main() {
	db := app.MustConnectDB()
	repo := repository.NewUserRepository(db)
	h := handler.NewUserHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.GetUsers)
	mux.HandleFunc("/common-friends", h.GetCommonFriends)
	mux.HandleFunc("/users/soft-delete", h.SoftDeleteUser)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}