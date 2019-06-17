package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "username or password not found", http.StatusBadRequest)
        return
	}
	var account Account
	//Create the account
	account, err = CreateMainAccount(Database, credentials.Username, []byte(credentials.Password))
	if err != nil {
		//Respond to errors
		DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusConflict)
        return
	}
    jwt, err := GenToken(account.Username,account.TokValue)
	if err != nil {
        DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
        return
    }
    response := map[string]string{"jwt":jwt}
    var encodedResponse []byte
    encodedResponse, err = json.Marshal(response)
	if err != nil {
        DualDebug(fmt.Sprintf("Error reading body: %v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
        return
    }
    w.Write(encodedResponse)
	DualDebug(fmt.Sprintf("%s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path))
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
	}
	router := mux.NewRouter()
	DualInfo("Loading Routes")
	//Loop over the list of Routes and add them to the router
	for _, r := range TheRoutes {
		router.HandleFunc(r.Path, r.Handler).Methods(r.Method)
	}
	return router
}
