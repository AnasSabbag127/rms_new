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
	// "github.com/gorilla/mux"
	// "strconv"
)

func CreateDishesHandler(w http.ResponseWriter, r *http.Request) {
	var newDish model.InputDishes
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newDish); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()

	if err != nil {
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	sql := `INSERT INTO dishes(restraunt_id,created_by,name,price) VALUES($1,$2,$3,$4)`
	_, err = db.Exec(sql, newDish.RestaurantId, newDish.CreatedBy, newDish.DishName, newDish.Price)

	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Dish created successfully",
		"dish":    newDish,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

}

func UpdateDishHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	var updateDish model.Dishes
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateDish); err != nil {
		log.Println("Invalid Update Dish body : ", err)
		http.Error(w, "Invalid Updated Dish body: ", http.StatusBadRequest)
		return
	}
	pathParam := mux.Vars(r)
	dishID, err := strconv.Atoi(pathParam["dishId"])
	if err != nil {
		log.Println("Invalid path variable dish ID : ", err)
		http.Error(w, "Invalid path variable dish ID : ", http.StatusBadRequest)
		return
	}
	if dishID != updateDish.Id {
		log.Println("dishId not match with updateDish body ")
		http.Error(w, "Updated Dish:Invalid input data ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()

	if err != nil {
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	sql := `UPDATE  dishes SET restraunt_id=$1,created_by=$2,name=$3,price=$4 WHERE id =$5`

	_, err = db.Exec(sql, updateDish.RestaurantId, updateDish.CreatedBy, updateDish.DishName, updateDish.Price, dishID)

	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		http.Error(w, "DATABASE ERROR: ", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Dish Updated successfully",
		"dish":    updateDish,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
}

func DeleteDishHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	pathParam := mux.Vars(r)

	dishId := r.URL.Query().Get("dishId")
	dishID, err := strconv.Atoi(dishId)
	if err != nil {
		log.Println("invalid query param dishId", err)
		http.Error(w, "Invalid query param dishId ", http.StatusBadRequest)
		return
	}
	//restId from path variable
	restId, err := strconv.Atoi(pathParam["restId"])
	if err != nil {
		log.Println("invalid path variable ", err)
		http.Error(w, "Invalid Path variable", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}

	sql := `UPDATE dishes SET is_deleted = TRUE WHERE restraunt_id = $1 AND id = $2`
	_, err = db.Exec(sql, restId, dishID)
	if err != nil {
		log.Println("SQL Error: ", err)
		http.Error(w, "SQL Error:", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "restaurant dish successfully soft deleted",
		"restId":  restId,
		"dishId":  dishId,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error: ", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
}
func GetAllDishHandler(w http.ResponseWriter, r *http.Request) {
	if middlewares.CheckForTokenValidation(w, r) == false {
		log.Println("Invalid token ")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	dishes := make([]model.Dishes, 0)
	pathParam := mux.Vars(r)
	restId, err := strconv.Atoi(pathParam["restId"])
	if err != nil {
		log.Println("Invalid path variable restId ", err)
		http.Error(w, "Invalid restId", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}
	sql := `SELECT id,name,price,created_by,restraunt_id FROM dishes WHERE restraunt_id =$1 AND is_deleted = FALSE`

	rows, err := db.Query(sql, restId)
	if err != nil {
		log.Println("DB Conn Err: ", err)
		http.Error(w, "DB Conn Err: ", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var dis model.Dishes
		err = rows.Scan(&dis.Id, &dis.DishName, &dis.Price, &dis.CreatedBy, &dis.RestaurantId)
		if err != nil {
			log.Println("SQL rows scan Err: ", err)
			http.Error(w, "SQL ROWs Scan Err: ", http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, dis)
	}
	response := map[string]interface{}{
		"message":       "List of dishes of a restaurant ",
		"Restaurant ID": restId,
		"Dishes":        dishes,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error: ", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

}
