package api

import (
	"github.com/gorilla/mux"
	"net/http"
	// "strconv"
)

//admin routes
func AdminRoutes(router *mux.Router) *mux.Router {
	//can't create direct sub-admin from users body

	//admin operations for sub-admin : creates routes for -> update sub-admin -> get all sub-admin
	router.Handle("/create-sub-admin/{userId}", http.HandlerFunc(CreateSubAdminHandler)).Methods("POST")
	router.Handle("/update-sub-admin/{userId}", http.HandlerFunc(UpdateSubAdminHandler)).Methods("PUT")
	router.Handle("/delete-sub-admin", http.HandlerFunc(DeleteSubAdminHandler)).Methods("DELETE")
	//router.Handle("/all-sub-admins", http.HandlerFunc(GetAllSubAdminHandler)).Methods("GET")

	//Admin: operations for users.
	router.Handle("/create-user", http.HandlerFunc(CreateUserHandler)).Methods("POST")
	router.Handle("/update user/{userId}", http.HandlerFunc(UpdateUserHandler)).Methods("PUT")
	router.Handle("/delete-user", http.HandlerFunc(DeleteUserHandler)).Methods("DELETE")
	router.Handle("/all-users", http.HandlerFunc(GetAllUsersHandler)).Methods("GET")

	//Admin: operations for restaurants.
	router.Handle("/create-restaurant", http.HandlerFunc(CreateRestaurantHandler)).Methods("POST")
	router.Handle("/update-restaurant/{restId}", http.HandlerFunc(UpdateRestaurantHandler)).Methods("PUT")
	router.Handle("/delete-restaurant", http.HandlerFunc(DeleteRestaurantHandler)).Methods("DELETE")
	router.Handle("/all-restaurants", http.HandlerFunc(GetAllRestaurantHandler)).Methods("GET")

	//Admin: operations for restaurant dishes.
	router.Handle("/create-dishes", http.HandlerFunc(CreateDishesHandler)).Methods("POST")
	router.Handle("/update-dish/{dishId}", http.HandlerFunc(UpdateDishHandler)).Methods("PUT")
	router.Handle("/rest/{restId}/delete-dish", http.HandlerFunc(DeleteDishHandler)).Methods("DELETE")
	router.Handle("/rest/{restId}/all-dish", http.HandlerFunc(GetAllDishHandler)).Methods("GET")

	return router
}
