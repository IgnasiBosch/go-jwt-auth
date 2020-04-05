package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	privKeyPath = "keys/jwt_key.rsa"
	pubKeyPath  = "keys/jwt_key.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func initKeys() {
	var err error
	signKeyByte, err := ioutil.ReadFile(privKeyPath)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyByte)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}

	verifyKeyByte, err := ioutil.ReadFile(pubKeyPath)
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyByte)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {
	//PUBLIC ENDPOINTS
	http.HandleFunc("/login", LoginHandler)

	//PROTECTED ENDPOINTS
	http.Handle("/resource/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectHandler)),
	))

	log.Println("Now listening...")
	_ = http.ListenAndServe(":9090", nil)
}

func main() {
	initKeys()
	StartServer()
}

//////////

func ProtectHandler(w http.ResponseWriter, r *http.Request) {
	resp := Response{"Gained access to protected resource"}
	JsonResponse(resp, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u UserCredentials

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "Error in request")
	}

	//validate user credentials
	if strings.ToLower(u.Email) != "test@example.com" {
		if u.Password != "1234" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			_, _ = fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	claims := jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().UTC().Add(time.Minute * 20).Unix(),
		"sub": struct {
			Email string `json:"email"`
			Role string `json:"role"`
		}{u.Email, "Member"},
	}

	//crete a rsa 256 signer
	signer := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	tokenString, err := signer.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "Error while signing the token")
	}

	//create a token instance using the token string
	response := Token{tokenString}
	JsonResponse(response, w)
}

//AUTH TOKEN VALIDATION
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//validate token
	tokenString := ExtractToken(r)
	token, err := ParseToken(tokenString)

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func ExtractToken(r *http.Request) string {
	headerAuth := r.Header.Get("Authorization")
	authValue := strings.Split(headerAuth, " ")
	if len(authValue) != 2 || authValue[0] != "Bearer" {
		return ""
	}

	return authValue[1]
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexppected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
}

// HELPER FUNCTIONS
func JsonResponse(r interface{}, w http.ResponseWriter) {
	jsonRes, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonRes)
}
