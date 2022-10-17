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
MIIEowIBAAKCAQEAuMVb3lrWKlmGIzTaJqtVJi2rPIy7/BkRKlAZ7Q1u0VlyOhzI
cXq6zAGsh31uWeJBJFKrZdwA6b2LD5vscnuilHi0nfQZA5l+meQT9LJ7STnfJ7f8
1CniBIhj5g6dOva9o/ljrLCmRSE4MjJRl3LkydvrHNokzicOAvieq4BYgHVJ2DC7
r7cSWrHeIiBIEBb1zAghc1OxtkFkxumva2gmywq0zB0VfzsYrtSpWT77qSNA+UEH
H+UCahefkun1uaBcOrEIPUh/j95N3rFTAdIAzFvDqLvLdQbmm/NT78O7izSGvJ1K
4rM2fSiecALSROZqOeJy3pf2v6Caqlp2gdDHUQIDAQABAoIBAEeOBrnhq7bS3KOd
wC3hhCQ442ubhOFoQ8GDK8clwJjKbvYaV3W69cQzkcEWzjl46YlLipzzyla61LPC
ypq7Tob5B9lzwowmUWT/csr8o8oD42vaUMtJPQJMX4OkfTdsfpyV5AfokTuMVdr6
qaZhFEEoLbEKud4sObzk023PUnbMTMyz2weQsmVLT9gm64K9GHzauoEDWGpr3Boz
+mEKo9rhySeK73qXqJFFATl0IB50GpdKondere/XofLQ2kqUB2fFEWIeRXawFd4u
6yE5nxjhIvn9pqMhW48vPAYtwTUhY8OxpCwk9U+awBYMvkCnFlrO4TwCNFqE8j3/
r6YN5VECgYEA2c0GqQVfwbcgB+8pLmd8H7dq+HqHToTA00pgNAsl/LRFUhkIVyY1
YUJkr7nJ3ogzqyjNBz0K4JeaMbMxJRwbxu4SUBgmvc+3h0oa9ibTKtBH4KWAqwCC
3TLQFoMXJCf+YJOL8ZNeqQz4rEpCVvJj86p9F4iTRhl5bd7mYqiEjusCgYEA2S1U
19iAIEiSETtUWEzj/iAKOUZJKl4Kz6PVfUMT1lqodOHsNHexzcQs+W/LW/Z9Jr7u
JtdC0ZZl2UCQ8PBTirsA/gimisnmVhx3mvKXOOO/Cf8RpN1/faTnQJj/Tu+h3aV7
n53Av2wSQPoWnnNB9W515gTYqwiV//CX44lGS7MCgYEAmjSoj4kniB8hBZ0WOi24
2zfg+/a80CH76F1TieWOysHUBtGEbze1OZxpb2WKgQ1MD9Y+e+6DQgr0eFXX6N9i
51DuFFlVLLThy17zge5xOnHnQi3L0Mb24Kg2XooIG2hZmYU94xelQOnXMx0MpUTO
8dl24e+n3kzxBZJ46cdIu2sCgYBzz5yizazlik16Ku07eSVLasKI8FYr5aJWP8Ok
3JRDhmy2h5NyFzIVzDs/eMI09Cig9MgCpl/XbCA7zhZ8pWunWzmYPfyxniDaYqvV
UPAbQjepmP9Lr2JBGiLHa88ZxOfITmqyH2mdqn/BbpuJO2U8//6W/pab/iQfK6mT
iKyXyQKBgGb49CT+zi1VCa4jh0EH78odTwRb/PleAw25CLVSj35MqQCEA2/jU92g
P8iTDU5mDNEaQMvMYN1fKiAPm6hiKp9/5Q02zQWNVzOywgc0L41Yb9K/SqxGwITa
dFwwxhKgBqzgCeIpvKLCGrhMztzPjUL9o4MSdNa92vajcilr8ld6
-----END RSA PRIVATE KEY-----`

	token := jwt.New(jwt.SigningMethodRS256)

	claims := jwt.MapClaims{}

	claims["exp"] = "1666007564"                           //time.Now().Add(time.Minute * 1080).Unix()
	claims["iat"] = "1665989564"                           //time.Now().Add(time.Minute * 1080).Unix()
	claims["jti"] = "7c86d08d-f61e-4c59-ab93-859f7a4c7398" //time.Now().Add(time.Minute * 1080).Unix()
	claims["iss"] = "https://tpm-keycloak.boquar.tech/auth/realms/GS"
	claims["aud"] = []string{
		"realm-management",
		"account",
	}

	claims["sub"] = "bb0ad8de-2321-48cf-adb5-e213666e3b8d"
	claims["typ"] = "Bearer"
	claims["azp"] = "dtvadmin"
	claims["session_state"] = "9b772692-facf-4b9f-aeda-661693952218"
	claims["scope"] = "profile email"
	claims["email_verified"] = "false"
	claims["name"] = "Dmitry Tyulkin"
	claims["preferred_username"] = "tyulkin.d"
	claims["given_name"] = "Dmitry"
	claims["family_name"] = "Tyulkin"
	claims["email"] = "tyulkin.d@galileosky.ru"

	claims["realm_access"] = map[string][]string{
		"roles": {"offline_access", "admin", "uma_authorization"},
	}
	claims["acr"] = "1"
	claims["resource_access"] = map[string]map[string][]string{
		"realm-management": {
			"roles": {
				"view-identity-providers",
				"view-realm",
				"manage-identity-providers",
				"impersonation",
				"realm-admin",
				"create-client",
				"manage-users",
				"query-realms",
				"view-authorization",
				"query-clients",
				"query-users",
				"manage-events",
				"manage-realm",
				"view-events",
				"view-users",
				"view-clients",
				"manage-authorization",
				"manage-clients",
				"query-groups",
			},
		},
		"account": {
			"roles": {
				"manage-account",
				"manage-account-links",
				"view-profile",
			},
		},
	}

	claims["session_state"] = "878f3ca4-a6d8-45dd-a494-955fb575282f"

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
