package api

import (
	"github.com/gorilla/mux"
	"net/http"
	// "strconv"
)

//admin routes
func AdminRoutes(router *mux.Router) *mux.Router {

	router.Handle("/admin-create-users", http.HandlerFunc(CreateUserHandler)).Methods("POST")

	router.Handle("/create-sub-admin", http.HandlerFunc(CreateUserHandler)).Methods("POST")

	router.Handle("/create-restaurant", http.HandlerFunc(CreateRestaurantHandler)).Methods("POST")
	//router.Handle("/update-restaurant/{restId}", http.HandlerFunc(UpdateRestaurantHandler)).Methods("PUT")
	//router.Handle("/delete-restaurant/{restId}", http.HandlerFunc(DeleteRestaurantHandler)).Methods("DELETE")
	//router.Handle("/all-restaurants", http.HandlerFunc(GetAllRestaurantHandler)).Methods("GET")

	router.Handle("/create-dishes", http.HandlerFunc(CreateDishesHandler)).Methods("POST")
	//router.Handle("/update-dish/{dishId}", http.HandlerFunc(UpdateDishHandler)).Methods("PUT")
	//router.Handle("/delete-dish/{dishId}", http.HandlerFunc(DeleteDishHandler)).Methods("DELETE")
	//router.Handle("/rest/{restId}/all-dish", http.HandlerFunc(GetAllDishHandler)).Methods("GET")

	return router
}
