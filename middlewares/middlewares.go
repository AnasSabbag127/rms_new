package middlewares
import (
	"net/http"
)
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Access-Controll-Allow-Origin","*")
		w.Header().Set("Access-Controll-Allow-Headers","*")
		w.Header().Set("Access-Controll-Allow-Methods","*")
		if r.Method == "OPTIONS"{
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("content-Type","application/json")
		next.ServeHTTP(w,r)
	})
}