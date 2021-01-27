package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Name struct {
	FirstName string
	LastName  string
}

type Message struct {
	FirstKey    string
	SecondKey   string
	Name        Name
	PhoneNumber string
	ICQ         string
	LastKey     int64
}

// var M []byte

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

var Log string
var Pass string

func main() {

	fmt.Println("Логин")
	fmt.Scanf("%s\n", &Log)

	fmt.Println("Пароль")
	fmt.Scanf("%s\n", &Pass)

	r := mux.NewRouter()

	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

var mySigningKey = []byte("secret")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	t := time.Now()
	login := r.FormValue("login")
	password := r.FormValue("password")

	data := r.FormValue("data")
	dataanswer := data + "1"

	fmt.Println("User login = ", login, "; Server login = ", Log, "; \nUser password = ", password, "; Server password = ", Pass)

	autorizationok := Log == login && Pass == password
	fmt.Println("autorizationok = ", autorizationok)

	claims["admin permissions?"] = "maybe"
	claims["login"] = &Log

	claims["Data answer is"] = dataanswer
	claims["Token request at"] = t
	claims["ATTENTION!"] = "Привет, Макс :)"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
	}

	tokenFprint := []byte(tokenString)

	if autorizationok {
		fmt.Fprint(w, fmt.Sprintf("Token request at [%s]\nUser:\nLogin: '%s'\nPassword: '%s'\nData answer is: %s\n", t.Format(time.RFC3339), login, password, dataanswer))
		fmt.Fprint(w, fmt.Sprintf("Token: %s", tokenFprint))
	} else {

		fmt.Fprint(w, " access denied ")
	}

})

var m = Message{"World", "Hello", Name{"Dmitry", "Victorovich"}, "79082706690", "393181839", 211}

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	payload, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
