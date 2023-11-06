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

//UserRoutes
//func UserRoutes(router *mux.Router) *mux.Router {
//	router.Handle("/get-all-restaurant", http.HandlerFunc(GetAllRestaurantHandler)).Methods("GET")
//	//router.Handle("/get-all-dishes-of restaurant", http.HandlerFunc(GetAllDishesHandler)).Methods("GET")
//	return router
//}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	var newUser model.InputUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	hashPsw, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password hashing error : ", err)
		http.Error(w, "password hashing error ", http.StatusInternalServerError)
		return
	}
	hashPswStr := string(hashPsw)

	db, err := database.ConnectToDB()
	if err != nil {
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	sql := `INSERT INTO usersNew(name,email,address,password,role_id,created_by) VALUES($1,$2,$3,$4,$5,$6)`
	_, err = db.Exec(sql, newUser.Name, newUser.Email, newUser.Address, hashPswStr, newUser.RoleId, newUser.CreatedBy)

	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "user created successfully",
		"user":    newUser,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	//get userId from path variable
	pathParam := mux.Vars(r)
	userId, err := strconv.Atoi(pathParam["userId"])
	if err != nil {
		log.Println("Invalid Path variable userId ")
		http.Error(w, "Invalid Path variable userId ", http.StatusBadRequest)
		return
	}
	var updateUser model.InputUser
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&updateUser); err != nil {
		log.Println("Invalid update user body : ", err)
		http.Error(w, "Invalid update user body : ", http.StatusBadRequest)
		return
	}
	//password := updateUser.Password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password hashing error for update user ", err)
		http.Error(w, "password hashing error for update user ", http.StatusInternalServerError)
		return
	}
	hashPassStr := string(hashPassword)
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Error: ", err)
		http.Error(w, "DB Connection Error ", http.StatusInternalServerError)
		return
	}

	SQL := `UPDATE usersNew  
                SET name = $1,
                    email = $2,
                    password = $3,
                    address = $4,
                    role_id = $5,
                    created_by =$6
				WHERE id = $7 
					AND is_deleted = FALSE;`

	_, err = db.Exec(SQL, updateUser.Name, updateUser.Email, hashPassStr, updateUser.Address, updateUser.RoleId, updateUser.CreatedBy, userId)
	if err != nil {
		log.Println("failed to execute query: ", err)
		http.Error(w, "failed to execute query: ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "user updated successfully",
		"user":    updateUser,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(" json Marshaling Error : ", err)
		http.Error(w, "json marshaling error ", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseJson)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	//get user id from query param
	userID := r.URL.Query().Get("userId")
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Error: ", err)
		http.Error(w, "DB Connection Error ", http.StatusInternalServerError)
		return
	}
	SQL := `UPDATE usersNew SET is_deleted = TRUE WHERE id =$1;`
	_, err = db.Exec(SQL, userID)
	if err != nil {
		log.Println("failed to execute query : ", err)
		http.Error(w, "failed to execute query ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "user successfully soft-deleted",
		"userId":  userID,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("json marshaling error: ", err)
		http.Error(w, "json marshaling error ", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseJson)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Error: ", err)
		http.Error(w, "DB Connection Error ", http.StatusInternalServerError)
		return
	}
	var users []model.OutputUser
	SQL := `SELECT id,
       			name,
       			email,
       			address,
       			role_id,
       			created_by 
			FROM usersNew 
            WHERE is_deleted = FALSE;`
	rows, err := db.Query(SQL)
	if err != nil {
		log.Println("failed to execute query :", err)
		http.Error(w, "failed to execute query ", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var user model.OutputUser
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.RoleId, &user.CreatedBy)
		if err != nil {
			log.Println("rows scan error : ", err)
			http.Error(w, "rows scan error ", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	response := map[string]interface{}{
		"message": "list of all users ",
		"users":   users,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("json marshaling error: ", err)
		http.Error(w, "json marshaling error: ", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseJson)
}
