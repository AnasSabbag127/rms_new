package api

import (
	"awesomeProject/database"
	"awesomeProject/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

//UserRoutes
func UserRoutes(router *mux.Router) *mux.Router {
	router.Handle("/get-all-restraunt", http.HandlerFunc(GetAllRestrauntHandler)).Methods("GET")
	router.Handle("/get-all-dishes", http.HandlerFunc(GetAllDishesHandler)).Methods("GET")
	return router
}

func GetAllRestrauntHandler(w http.ResponseWriter, r *http.Request) {

}
func GetAllDishesHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
