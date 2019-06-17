package main
import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
)
type Route struct {
    Name string
    Method string
    Path string
    Handler http.HandlerFunc
}

//Handler functions
func Placeholder (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Placeholder"))
    DualInfo(fmt.Sprintf("%s request from %s to %s", r.Method,r.RemoteAddr,r.URL.Path))
}

//Use Route structs to construct all the necessary routes
func InitRouter () *mux.Router {
// List of all routes for API
TheRoutes := []Route{
    Route {
        "Index",
        "Get",
        "/",
        Placeholder,
    },
    Route {
        "Register",
        "Post",
        "/register",
        Placeholder,
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

