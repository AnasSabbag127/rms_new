package routes

import (
	"awesomeProject/api"
	"awesomeProject/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRoutes() http.Handler {
	r := mux.NewRouter()
	r.Handle("/admin", api.AdminRoutes(r))
	r.Handle("/sub-admin", api.SubAdminRoutes(r))
	r.Handle("/user", api.UserRoutes(r))
	return middlewares.EnableCORS(r)
}
