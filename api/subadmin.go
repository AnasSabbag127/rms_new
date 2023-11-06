package api

import (
	"awesomeProject/database"
	"awesomeProject/middlewares"
	"awesomeProject/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

//subadmin Routes
func SubAdminRoutes(router *mux.Router) *mux.Router {

	//sub admin : user operations
	router.Handle("/create-user", http.HandlerFunc(CreateUserHandler)).Methods("POST")
	router.Handle("/update user/{userId}", http.HandlerFunc(UpdateUserHandler)).Methods("PUT")
	router.Handle("/delete-user", http.HandlerFunc(DeleteUserHandler)).Methods("DELETE")
	router.Handle("/all-users", http.HandlerFunc(GetAllUsersHandler)).Methods("GET")

	//sub admin : restaurant operations
	router.Handle("/create-restaurant", http.HandlerFunc(CreateRestaurantHandler)).Methods("POST")
	router.Handle("/update-restaurant/{restId}", http.HandlerFunc(UpdateRestaurantHandler)).Methods("PUT")
	router.Handle("/delete-restaurant", http.HandlerFunc(DeleteRestaurantHandler)).Methods("DELETE")
	router.Handle("/all-restaurants", http.HandlerFunc(GetAllRestaurantHandler)).Methods("GET")

	//sub admin : restaurant dishes operations
	router.Handle("/create-dishes", http.HandlerFunc(CreateDishesHandler)).Methods("POST")
	router.Handle("/update-dish/{dishId}", http.HandlerFunc(UpdateDishHandler)).Methods("PUT")
	router.Handle("/rest/{restId}/delete-dish", http.HandlerFunc(DeleteDishHandler)).Methods("DELETE")
	router.Handle("/rest/{restId}/all-dish", http.HandlerFunc(GetAllDishHandler)).Methods("GET")

	return router
}

func CreateSubAdminHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	//get user id for creating sub-admin from the path variable
	userId := mux.Vars(r)
	userID, err := strconv.Atoi(userId["userId"])
	if err != nil {
		log.Println("invalid userId : ", err)
		http.Error(w, "invalid userId ", http.StatusBadRequest)
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Failed to connect to Database: ", err)
		http.Error(w, "Failed to connect to Database ", http.StatusInternalServerError)
		return
	}
	subAdminRoleId := 2
	// for creating sub-admin update role_id if user is not exist.
	SQL := `UPDATE usersNew SET role_id = $1 WHERE id = $2 AND is_deleted = false;`
	_, err = db.Exec(SQL, subAdminRoleId, userID)
	if err != nil {
		log.Println("failed to execute query: ", err)
		http.Error(w, "failed to execute query", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":      "sub-admin created successfully",
		"sub-admin id": userID,
	}
	responseJson, err := json.Marshal(response)
	_, err = w.Write(responseJson)

}

func UpdateSubAdminHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("invalid token: ")
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	//
	pathParam := mux.Vars(r)
	userId, err := strconv.Atoi(pathParam["userId"])
	if err != nil {
		log.Println("invalid path variable : ", err)
		http.Error(w, "invalid path variable userId : ", http.StatusBadRequest)
		return
	}
	//sub-admin: update body
	var updateBody model.InputUser
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&updateBody); err != nil {
		log.Println("invalid sub-admin update body : ", err)
		http.Error(w, "invalid sub-admin update body ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("failed to connect DB : ", err)
		http.Error(w, "failed to connect database", http.StatusInternalServerError)
		return
	}
	passwordStr := updateBody.Password

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
	if err != nil {
		log.Println(" password hashing error: ", err)
		http.Error(w, "password hashing error ", http.StatusInternalServerError)
		return
	}
	hashPasswordStr := string(hashPassword)
	//update sub-admin query
	SQL := `UPDATE usersNew SET
				name = $1,
				email = $2,
				address = $3,
				password = $4,
				role_id = $5
			WHERE id = $6 AND is_deleted = false AND role_id = 2;
			`

	_, err = db.Exec(SQL, updateBody.Name, updateBody.Email, updateBody.Address, hashPasswordStr,
		updateBody.RoleId, userId)

	if err != nil {
		log.Println("failed to execute query : ", err)
		http.Error(w, "failed to execute query ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":   "sub-admin updated successfully",
		"sub-admin": updateBody,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("json marshaling error ", err)
		http.Error(w, "json marshaling error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseJson)
}
func DeleteSubAdminHandler(w http.ResponseWriter, r *http.Request) {
	//this handler function only remove or updated from sub-admin.
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("invalid token: ")
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	//get userId from query param
	subAdminId := r.URL.Query().Get("userId")
	subAdminID, err := strconv.Atoi(subAdminId)
	if err != nil {
		log.Println("path variable : invalid sub admin id ", err)
		http.Error(w, "invalid user id ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("failed to connect DB : ", err)
		http.Error(w, "failed to connect database", http.StatusInternalServerError)
		return
	}
	userRole := 3
	SQL := `UPDATE usersNew SET role_id = $1 WHERE id = $2  AND is_deleted = false`
	_, err = db.Exec(SQL, userRole, subAdminID)
	if err != nil {
		log.Println("failed to execute query : ", err)
		http.Error(w, "failed to execute query ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":   "sub-admin removed successfully",
		"sub-admin": subAdminID,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("json marshaling error ", err)
		http.Error(w, "json marshaling error ", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseJson)

}

func GetAllSubAdminHandler(w http.ResponseWriter, r *http.Request) {

}

//func GetAllUserCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
//
//	var subAdmin model.User
//	// get sub-admin info  from path variable later
//	var users []model.OutputUser
//	db, err := database.ConnectToDB()
//	if err != nil {
//		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
//		return
//	}
//	//get all user for sub-admin it will get only those users which is created by him
//	SQL := `SELECT id,name,email,address,role_id FROM usersNew WHERE created_by =$1`
//
//	rows, err := db.Query(SQL, subAdmin.Id)
//	if err != nil {
//		log.Println("DATABASE ERROR: ", err)
//		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
//		return
//	}
//	for rows.Next() {
//		var user model.OutputUser
//		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.RoleId)
//		if err != nil {
//			log.Println("rows scan error: ", err)
//			http.Error(w, "sub-admin get users :", http.StatusInternalServerError)
//			return
//		}
//		users = append(users, user)
//	}
//
//	response := map[string]interface{}{
//		"message": "users fetched successfully",
//		"user":    users,
//	}
//	jsonResponse, err := json.Marshal(response)
//	if err != nil {
//		log.Println("json Marshaling error:  ", err)
//		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
//		return
//	}
//	_, _ = w.Write(jsonResponse)
//
//}
//
//func GetAllRestaurantsCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
//	var subAdminID int
//	var restaurants []model.Restaurant
//	db, err := database.ConnectToDB()
//	if err != nil {
//		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
//		return
//	}
//	//get all restaurants: for sub-admin it will get only those restaurants which is created by him
//	SQL := `SELECT id,name,address,stars,created_by FROM restraunts WHERE created_by =$1`
//
//	rows, err := db.Query(SQL, subAdminID)
//	if err != nil {
//		log.Println("DATABASE ERROR: ", err)
//		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
//		return
//	}
//	for rows.Next() {
//		var rest model.Restaurant
//		err = rows.Scan(&rest.Id, &rest.Name, &rest.Address, &rest.Stars, &rest.CreatedBy)
//		if err != nil {
//			log.Println("rows scan error: ", err)
//			http.Error(w, "sub-admin get all restaurants :", http.StatusInternalServerError)
//			return
//		}
//		restaurants = append(restaurants, rest)
//	}
//	response := map[string]interface{}{
//		"message":     "restaurants fetched  successfully",
//		"restaurants": restaurants,
//	}
//	jsonResponse, err := json.Marshal(response)
//	if err != nil {
//		log.Println("json Marshaling error:  ", err)
//		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
//		return
//	}
//	_, _ = w.Write(jsonResponse)
//}
//
//func GetAllDishesOfRestaurantCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
//	var restID int
//	var subAdminID int
//	var dishes []model.Dishes
//	db, err := database.ConnectToDB()
//	if err != nil {
//		log.Println("failed to connect to DB: ", err)
//		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
//		return
//	}
//	SQL := `SELECT id,name,price,created_by,restraunt_id FROM dishes WHERE created_by = $1 AND restraunt_id = $2`
//	rows, err := db.Query(SQL, subAdminID, restID)
//	for rows.Next() {
//		var dish model.Dishes
//		err = rows.Scan(&dish.Id, &dish.DishName, &dish.Price, &dish.CreatedBy, &dish.RestaurantId)
//		if err != nil {
//			log.Println("rows scan error: ", err)
//			http.Error(w, "sub-admin get all dishes of restaurant Error: ", http.StatusInternalServerError)
//			return
//		}
//		dishes = append(dishes, dish)
//	}
//	response := map[string]interface{}{
//		"message": "dishes fetched  successfully",
//		"dishes":  dishes,
//	}
//	jsonResponse, err := json.Marshal(response)
//	if err != nil {
//		log.Println("json Marshaling error:  ", err)
//		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
//		return
//	}
//	_, _ = w.Write(jsonResponse)
//}
