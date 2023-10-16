package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

//subadmin Routes
func SubAdminRoutes(router *mux.Router) *mux.Router {
	router.Handle("/create-users", http.HandlerFunc(CreateUserHandler)).Methods("POST")
	router.Handle("/restaurant/{restId}/create-dish", http.HandlerFunc(CreateDishesHandler)).Methods("POST")
	router.Handle("/create-restaurant", http.HandlerFunc(CreateRestaurantHandler)).Methods("POST")
	return router
}
