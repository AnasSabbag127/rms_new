package middlewares

import (
	"awesomeProject/database"
	"awesomeProject/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("secret-key")

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credential
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println("Error is : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("DB conn error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var userAcc model.User
	err = db.QueryRow("SELECT password FROM usersNew WHERE name = $1", credentials.Username).Scan(&userAcc.Password)
	if err != nil {
		log.Println("user not found : ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// password hashing : apply later
	hashBytePassword := []byte(userAcc.Password)
	err = bcrypt.CompareHashAndPassword(hashBytePassword, []byte(credentials.Password))
	if err != nil {
		log.Println("Invalid Credentials: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			log.Println("cookie not found : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println("Error is: ", err)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !tkn.Valid {
		if errors.Is(err, jwt.ErrSignatureInvalid) || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("hello , %s", claims.Username)))
	w.WriteHeader(http.StatusOK)
	return
}

func CheckForTokenValidation(w http.ResponseWriter, r *http.Request) bool {
	//this will check for validation of token if its return true then api execute else not
	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			log.Println("cookie not found : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		log.Println("Error is: ", err)
		return false
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !tkn.Valid {
		if errors.Is(err, jwt.ErrSignatureInvalid) || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Controll-Allow-Origin", "*")
		w.Header().Set("Access-Controll-Allow-Headers", "*")
		w.Header().Set("Access-Controll-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
