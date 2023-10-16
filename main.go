package main

import (
	"awesomeProject/database"
	"awesomeProject/routes"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

func handler(_ http.ResponseWriter, r *http.Request) {

}
func main() {
	fmt.Println("jfjkf")
	_, err := database.ConnectToDB()
	if err != nil {
		fmt.Println("database not connected: Error: ", err)
		return
	}
	srv := routes.CreateRoutes()

	// http.HandleFunc("/",handler)
	// http.HandleFunc("/login",handler)
	// //admin routes
	// http.HandleFunc("/create-user",api.CreateUserHandler)
	// http.HandleFunc("/user/create-restraunt",api.CreateRestrauntHandler)
	// http.HandleFunc("/user/restraunt/create-dishes",api.CreateDishesHandler)

	// http.ListenAndServe("localhost:8000",nil)

	if err := http.ListenAndServe("localhost:8000", handlers.RecoveryHandler()(srv)); err != nil {
		fmt.Println("ListenAndServe Errors:", err)
		return
	}

}
