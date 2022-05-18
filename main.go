package main

import (
	"context"
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

	// redis
	My_redis "github.com/ev/myhttp3authjson/Cache"
	"github.com/go-redis/redis/v8"
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

type Gettokenanswerstruct struct {
	TokenRequestAt string
	User           string
	Login          string
	Password       string
	DataAnswer     string
	Token          string
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

var Log string
var Pass string

var ctx = context.Background()

func main() {
	red := My_redis.MyRedis{}
	red.Name = "abc"
	//redis start
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println(red)

	userid := "user#125"
	currentusertoken := "testtoken"

	node := rdb.Set(ctx, userid, currentusertoken, 0).Err()
	if node != nil {
		panic(node)
	}

	val, node := rdb.Get(ctx, userid).Result()
	if node == redis.Nil {
		fmt.Println("key1 does not exist")
	} else if node != nil {
		panic(node)
	} else {
		fmt.Println(userid, val)
	}

	val2, node := rdb.Get(ctx, "key2").Result()
	if node == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if node != nil {
		panic(node)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist

	//redis end

	Log = "root111"
	Pass = "1111"

	r := mux.NewRouter()

	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

var mySigningKey = []byte("secret")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	t := time.Now()
	login := r.FormValue("login")
	password := r.FormValue("password")

	data := r.FormValue("data")
	dataanswer := data + "1"

	fmt.Println("User login = ", login, "; Server login = ", Log, "; \nUser password = ", password, "; Server password = ", Pass)

	autorizationok := Log == login && Pass == password

	fmt.Println("autorizationok = ", autorizationok)

	newuserid := ""

	if autorizationok {
		newuserid = Log + Pass
	} else {
		newuserid = "false"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val3, newuserautorizationdata := rdb.Get(ctx, newuserid).Result()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin permissions?"] = "maybe"
	claims["login"] = &Log

	claims["Data answer is"] = dataanswer
	claims["Token request at"] = t
	claims["ATTENTION!"] = "Привет, Макс :)"
	claims["exp"] = time.Now().Add(time.Minute * 1080).Unix()
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatal(err)
	}

	var actualtokenstring string

	if newuserautorizationdata == redis.Nil {
		fmt.Println("newuserautorizationdata does not exist")
		actualtokenstring = tokenString

	} else if newuserautorizationdata != nil {
		panic(newuserautorizationdata)
	} else {
		fmt.Println(newuserid, val3)
		actualtokenstring = val3
	}

	currentuserautorizationdata := rdb.Set(ctx, newuserid, actualtokenstring, 0).Err()
	if currentuserautorizationdata != nil {
		panic(currentuserautorizationdata)
	}

	fmt.Println("NewToken = ", actualtokenstring)

	//tokenFprint := []byte(actualtokenstring)

	var Gettokenanswer = Gettokenanswerstruct{t.Format(time.RFC3339), newuserid, login, password, dataanswer, actualtokenstring}

	if autorizationok {
		// fmt.Fprint(w, fmt.Sprintf("Token request at [%s]\nUser:\nLogin: '%s'\nPassword: '%s'\nData answer is: %s\n", t.Format(time.RFC3339), login, password, dataanswer))
		// fmt.Fprint(w, fmt.Sprintf("Token: %s", tokenFprint))
		payload, _ := json.Marshal(Gettokenanswer)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
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
