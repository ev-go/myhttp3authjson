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
	My_redis "github.com/ev-go/myhttp3authjson/Cache"
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
	SigningMethod: jwt.SigningMethodRS256,
})

var Log string
var Pass string

var ctx = context.Background()

func main() {
	My_redis.Main()

	Log = "root22"
	Pass = "1"

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

	var key = `-----BEGIN RSA PRIVATE KEY-----
MIIBUwIBADANBgkqhkiG9w0BAQEFAASCAT0wggE5AgEAAkEAqaaVjLmQXez9pc0/0LjUgjScrwap5XYKVdU1P7MP0cFK46k+g4jw2tZWO/tHgpHCPqcBxazvXk6quSP5x/RH4QIDAQABAkAFLWBC+ya8tw7GkIzyAbH6h0CA1uM4bXHDTM5jM0O4UvEGFXKOwfitm4r5YODBYpxbCtGBq2WNbZnOkl/VmExBAiEA31Q0RzjYx+9iIvoXJyzTE5joYJkMrLOj2DZWXQvDWHsCIQDCeBm/dR3M0EokSp0yei2QOlSTSbHEm+1679YCHXRIUwIgDe8aB+uTEv5rmBUUhrw0Oz/KF+TUtp3Ktj3Fq66FNKsCIAzjFpm4ciQbfX5QL4Cj1hcjtm0YSh6EUsV91UnIl+cfAiAc9fQWN/jg9BzGpnW70HAbbQJT7+nkK4H0dCKriEqV3A==
-----END RSA PRIVATE KEY-----`

	token := jwt.New(jwt.SigningMethodRS256)

	claims := jwt.MapClaims{}

	claims["exp"] = "1663927554"                           //time.Now().Add(time.Minute * 1080).Unix()
	claims["iat"] = "1663927254"                           //time.Now().Add(time.Minute * 1080).Unix()
	claims["jti"] = "14ccf931-544d-4f49-82ce-2e6a3ca0c1ec" //time.Now().Add(time.Minute * 1080).Unix()
	claims["iss"] = "https://tpm-keycloak.boquar.tech/auth/realms/GS"
	//claims["aud"] = []string{
	//	"realm-management",
	//	"account",
	//}

	claims["aud"] = "account"

	claims["sub"] = "c5074793-9d82-478c-9853-125c04bdb626"
	claims["typ"] = "Bearer"
	claims["azp"] = "ng-frontend"
	claims["session_state"] = "77f26529-fbd5-4cfb-a872-561c1267c890"
	claims["acr"] = "1"
	claims["realm_access"] = map[string][]string{
		"roles": {"offline_access", "admin", "uma_authorization", "support"},
	}

	claims["resource_access"] = map[string]map[string][]string{
		"account": {
			"roles": {
				"manage-account",
				"manage-account-links",
				"view-profile",
			},
		},
	}

	claims["scope"] = "profile email"
	claims["email_verified"] = "false"
	claims["group_uuids"] = []string{
		"0da3b22f-ec3f-4383-bc25-480b6dcb82a1",
		"5ebc3080-1023-47dd-8927-d704ef377e9d",
	}
	claims["customer_uuid"] = "0da3b22f-ec3f-4383-bc25-480b6dcb82a1"
	claims["name"] = "Инженер Техподдержки"
	claims["preferred_username"] = "support"
	claims["given_name"] = "Инженер"
	claims["family_name"] = "Техподдержки"
	claims["email"] = "example@test.ru"

	claims["session_state"] = "878f3ca4-a6d8-45dd-a494-955fb575282f"

	claims["group"] = []string{
		"/Galileosky",
		"/Galileosky/technical support",
	}

	//token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header = map[string]interface{}{"kid": "JYyQHHvNTMBntTKc8-m5kooVWLk8hXKWDVrc56bw15E",
		"alg": "HS256",
		"typ": "JWT",
	}

	token.Claims = claims

	//signBytes := []byte("MIIBOgIBAAJBAKj34GkxFhD90vcNLYLInFEX6Ppy1tPf9Cnzj4p4WGeKLs1Pt8QuKUpRKfFLfRYC9AIKjbJTWit+CqvjWYzvQwECAwEAAQJAIJLixBy2qpFoS4DSmoEmo3qGy0t6z09AIJtH+5OeRV1be+N4cDYJKffGzDa88vQENZiRm0GRq6a+HPGQMd2kTQIhAKMSvzIBnni7ot/OSie2TmJLY4SwTQAevXysE2RbFDYdAiEBCUEaRQnMnbp79mxDXDf6AU0cN/RPBjb9qSHDcWZHGzUCIG2Es59z8ugGrDY+pxLQnwfotadxd+Uyv/Ow5T0q5gIJAiEAyS4RaI9YG8EWx/2w0T67ZUVAw8eOMB6BIUg0Xcu+3okCIBOs/5OiPgoTdSy7bcF9IGpSE8ZgGKzgYQVZeN97YE00")
	//signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	//if err != nil {
	//	log.Fatal(err)
	//}

	//tokenString, err := token.SignedString(key)

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal(err)
	}
	tokenString, err := token.SignedString(pkey)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("this is fresh token", tokenString)
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

	var Gettokenanswer = Gettokenanswerstruct{t.Format(time.RFC3339), newuserid, login, password, dataanswer, tokenString}

	if autorizationok {
		payload, _ := json.Marshal(Gettokenanswer)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
	} else {

		fmt.Fprint(w, " access denied ")
	}

})

var m = Message{
	"World",
	"Hello",
	Name{
		"Dmitry",
		"Victorovich"},
	"79082706690",
	"393181839",
	211}

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	payload, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
