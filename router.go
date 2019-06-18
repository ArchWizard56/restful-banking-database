package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
    "strings"
    "errors"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

//Handler functions
func Placeholder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder"))
	DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
}
func Register(w http.ResponseWriter, r *http.Request) {
	//Load credentials from request body
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "username or password not found", http.StatusBadRequest)
		return
	}
	//Create an account using the dataBaseInterface.go file, loading it into "account"
	var account Account
	//Create the account
	account, err = CreateMainAccount(Database, credentials.Username, []byte(credentials.Password))
	if err != nil {
		//Respond to errors
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusConflict)
		return
	}
	//Use the account details to create a JSON Web Token using the function in authorization.go
    var token string
    token, err = GetToken(Database, account.Username)
	if err != nil {
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	jwt, err := GenToken(account.Username, token)
	if err != nil {
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	//Encode the response as JSON and send it back to the client as a response
	response := map[string]string{"jwt": jwt}
	var encodedResponse []byte
	encodedResponse, err = json.Marshal(response)
	if err != nil {
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.Write(encodedResponse)
	//Debug output
	DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "username or password not found", http.StatusBadRequest)
		return
	}
	var ValidAccount bool
	ValidAccount, err = IsAccountValid(Database, credentials.Username, credentials.Password)
	if ValidAccount == true {
		var TokValue string
		TokValue, err = GetToken(Database, credentials.Username)
		jwt, err := GenToken(credentials.Username, TokValue)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		response := map[string]string{"jwt": jwt}
		var encodedResponse []byte
		encodedResponse, err = json.Marshal(response)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		w.Write(encodedResponse)
		DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
	} else {
		http.Error(w, "bad username or password", http.StatusUnauthorized)
	}
}
func LogOut(w http.ResponseWriter, r *http.Request) {
    if len(r.Header["Authorization"]) != 1 {
        err := errors.New("Bad Request")
        DualDebug(fmt.Sprintf("%v", err))
        http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
        return
    }
    authHeader := r.Header["Authorization"][0]
	Token := strings.Fields(authHeader)[1]
    ValidAccount, err := VerifyToken(Token)
	if err != nil {
		if err.Error() == "Invalid Token Value" || err.Error() == "Invalid Token" {
			DualDebug(fmt.Sprintf("%v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
			return
        }
        DualDebug(fmt.Sprintf("%v", err))
        http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
        return
	}
	claims, err := GetTokenClaims(Token)
	if err != nil {
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if ValidAccount == true {
		err = ChangeToken(Database, claims.Username)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		response := map[string]string{"result": "success"}
		var encodedResponse []byte
		encodedResponse, err = json.Marshal(response)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		w.Write(encodedResponse)
		DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
	} else {
		http.Error(w, "Bad AuthToken", http.StatusUnauthorized)
	}
}
func LoadAccounts(w http.ResponseWriter, r *http.Request) {
    if len(r.Header["Authorization"]) != 1 {
        err := errors.New("Bad Request")
        DualDebug(fmt.Sprintf("%v", err))
        http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
        return
    }
    authHeader := r.Header["Authorization"][0]
	Token := strings.Fields(authHeader)[1]
    ValidAccount, err := VerifyToken(Token)
	if err != nil {
		if err.Error() == "Invalid Token Value" || err.Error() == "Invalid Token" {
			DualDebug(fmt.Sprintf("%v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
			return
        }
        DualDebug(fmt.Sprintf("%v", err))
        http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
        return
	}
	claims, err := GetTokenClaims(Token)
	if err != nil {
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if ValidAccount == true {
        accounts, err:= GetAccounts(Database,claims.Username)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		var encodedResponse []byte
		encodedResponse, err = json.Marshal(accounts)
		if err != nil {
			DualDebug(fmt.Sprintf("Error reading body: %v", err))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		w.Write(encodedResponse)
		DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
	} else {
		http.Error(w, "Bad AuthToken", http.StatusUnauthorized)
	}
}

//Use Route structs to construct all the necessary routes
func InitRouter() *mux.Router {
	// List of all routes for API
	TheRoutes := []Route{
		Route{
			"Index",
			"Get",
			"/",
			Placeholder,
		},
		Route{
			"Register",
			"Post",
			"/register",
			Register,
		},
		Route{
			"SignIn",
			"Post",
			"/signin",
			SignIn,
		},
		Route{
			"Logout",
			"Post",
			"/logout",
			LogOut,
		},
		Route{
			"Accounts",
			"Get",
			"/accounts",
			LoadAccounts,
		},
	}
	router := mux.NewRouter()
	DualInfo("Loading Routes")
	//Loop over the list of Routes and add them to the router
	for _, r := range TheRoutes {
		router.HandleFunc(r.Path, r.Handler).Methods(r.Method)
	}
	return router
}
