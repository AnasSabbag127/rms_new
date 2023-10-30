package routes

import (
	"awesomeProject/api"
	"awesomeProject/middlewares"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "home handler working",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("json Marshaling error:  ", err)
		http.Error(w, "json marshaling Error", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResponse)
}

func CreateRoutes() http.Handler {
	r := mux.NewRouter()
	r.Handle("/home", http.HandlerFunc(middlewares.Home)).Methods("GET")
	r.Handle("/login", http.HandlerFunc(middlewares.Login)).Methods("POST")
	r.Handle("/", api.UserRoutes(r))
	r1 := r.PathPrefix("/").Subrouter()
	authRoutes(r1)
	return middlewares.EnableCORS(r)
}

func authRoutes(r *mux.Router) {
	r.Handle("/", api.AdminRoutes(r))
	r.Handle("/", api.SubAdminRoutes(r))
}
