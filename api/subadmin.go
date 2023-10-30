package api

import (
	"awesomeProject/database"
	"awesomeProject/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//subadmin Routes
func SubAdminRoutes(router *mux.Router) *mux.Router {
	router.Handle("/create-users", http.HandlerFunc(CreateUserHandler)).Methods("POST")

	router.Handle("/create-restaurant-by-sub-admin", http.HandlerFunc(CreateRestaurantHandler)).Methods("POST")
	router.Handle("/create-dish-by-sub-admin", http.HandlerFunc(CreateDishesHandler)).Methods("POST")

	//router.Handle("/get-all-users", http.HandlerFunc(GetAllUserCreatedBySubAdmin)).Methods("GET")
	//router.Handle("/{subAdminId}/get-all-restaurants", http.HandlerFunc(GetAllRestaurantsCreatedBySubAdmin)).Methods("GET")
	//router.Handle("/{subAdminId}/restaurant/{restId}/all-dishes", http.HandlerFunc(GetAllDishesOfRestaurantCreatedBySubAdmin)).Methods("GET")

	return router
}

func GetAllUserCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {

	var subAdmin model.User
	// get sub-admin info  from path variable later
	var users []model.OutputUser
	db, err := database.ConnectToDB()
	if err != nil {
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	//get all user for sub-admin it will get only those users which is created by him
	SQL := `SELECT id,name,email,address,role_id FROM usersNew WHERE created_by =$1`

	rows, err := db.Query(SQL, subAdmin.Id)
	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var user model.OutputUser
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.RoleId)
		if err != nil {
			log.Println("rows scan error: ", err)
			http.Error(w, "sub-admin get users :", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := map[string]interface{}{
		"message": "users fetched successfully",
		"user":    users,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResponse)

}

func GetAllRestaurantsCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
	var subAdminID int
	var restaurants []model.Restaurant
	db, err := database.ConnectToDB()
	if err != nil {
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	//get all restaurants: for sub-admin it will get only those restaurants which is created by him
	SQL := `SELECT id,name,address,stars,created_by FROM restraunts WHERE created_by =$1`

	rows, err := db.Query(SQL, subAdminID)
	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var rest model.Restaurant
		err = rows.Scan(&rest.Id, &rest.Name, &rest.Address, &rest.Stars, &rest.CreatedBy)
		if err != nil {
			log.Println("rows scan error: ", err)
			http.Error(w, "sub-admin get all restaurants :", http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, rest)
	}
	response := map[string]interface{}{
		"message":     "restaurants fetched  successfully",
		"restaurants": restaurants,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResponse)
}

func GetAllDishesOfRestaurantCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
	var restID int
	var subAdminID int
	var dishes []model.Dishes
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("failed to connect to DB: ", err)
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	SQL := `SELECT id,name,price,created_by,restraunt_id FROM dishes WHERE created_by = $1 AND restraunt_id = $2`
	rows, err := db.Query(SQL, subAdminID, restID)
	for rows.Next() {
		var dish model.Dishes
		err = rows.Scan(&dish.Id, &dish.DishName, &dish.Price, &dish.CreatedBy, &dish.RestaurantId)
		if err != nil {
			log.Println("rows scan error: ", err)
			http.Error(w, "sub-admin get all dishes of restaurant Error: ", http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, dish)
	}
	response := map[string]interface{}{
		"message": "dishes fetched  successfully",
		"dishes":  dishes,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResponse)
}
