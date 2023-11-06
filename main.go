package main

import (
	"awesomeProject/database"
	"awesomeProject/routes"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

func main() {
	fmt.Println("RMS")
	_, err := database.ConnectToDB()
	if err != nil {
		fmt.Println("database not connected: Error: ", err)
		return
	}
	srv := routes.CreateRoutes()

	if err := http.ListenAndServe("localhost:8000", handlers.RecoveryHandler()(srv)); err != nil {
		fmt.Println("ListenAndServe Errors:", err)
		return
	}

}
