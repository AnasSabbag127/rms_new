package api

import (
	"awesomeProject/database"
	"awesomeProject/middlewares"
	"awesomeProject/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// make restraunts api(CRUD) routes here

func CreateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	var newRest model.InputRestaurant
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newRest); err != nil {
		log.Println("Invalid Restaurant body : ", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	db, err := database.ConnectToDB()

	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	sql := `INSERT INTO restraunts(name,address,created_by,stars) VALUES($1,$2,$3,$4)`
	_, err = db.Exec(sql, newRest.Name, newRest.Address, newRest.CreatedBy, newRest.Stars)

	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":    "Restaurant created successfully",
		"Restaurant": newRest,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
}

func UpdateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	var updateRest model.Restaurant
	decode := json.NewDecoder(r.Body)

	if err := decode.Decode(&updateRest); err != nil {
		log.Println("Invalid UpdateRestaurant body : ", err)
		http.Error(w, "Invalid UpdateRestaurant body:", http.StatusBadRequest)
		return
	}
	pathParam := mux.Vars(r)
	restId, err := strconv.Atoi(pathParam["restId"])
	if err != nil {
		log.Println("Invalid path id {restId} : ", err)
		http.Error(w, "Invalid path variable : ", http.StatusBadRequest)
		return
	}
	if updateRest.Id != restId {
		log.Println("rest ID not match with updated rest-id : ", err)
		http.Error(w, "Invalid path variable : ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()

	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	sql := `UPDATE  restraunts SET name=$1,address=$2,created_by=$3,stars=$4 WHERE id = $5`
	_, err = db.Exec(sql, updateRest.Name, updateRest.Address, updateRest.CreatedBy, updateRest.Stars, updateRest.Id)
	if err != nil {
		log.Println("DB SQL Err: ", err)
		http.Error(w, "SQL Query Error ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":    "Restaurant Updated successfully",
		"Restaurant": updateRest,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

}

func DeleteRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	restId := r.URL.Query().Get("restId")
	restID, err := strconv.Atoi(restId)
	if err != nil {
		log.Println("invalid query param : ", err)
		http.Error(w, "Invalid query param", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}
	//apply soft deletion for deleting row of any table
	sql := `UPDATE restraunts SET is_deleted = TRUE WHERE id =$1`
	_, err = db.Exec(sql, restID)
	if err != nil {
		log.Println("SQL Error: ", err)
		http.Error(w, "SQL Error:", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "restaurant successfully deleted",
		"restId":  restId,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error: ", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
}

func GetAllRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	restaurants := make([]model.Restaurant, 0)
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}
	sql := `SELECT id,name,address,created_by,stars FROM restraunts WHERE is_deleted = FALSE ;`

	rows, err := db.Query(sql)
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var rest model.Restaurant
		err = rows.Scan(&rest.Id, &rest.Name, &rest.Address, &rest.CreatedBy, &rest.Stars)
		if err != nil {
			log.Println("SQL rows scan Err: ", err)
			http.Error(w, "SQL ROWs Scan Err: ", http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, rest)
	}
	response := map[string]interface{}{
		"message":     "lists of restaurants ",
		"Restaurants": restaurants,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error: ", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

}
