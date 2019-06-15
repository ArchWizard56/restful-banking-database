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
    fmt.Println(r.Method + " request from " + r.RemoteAddr + " to " + r.URL.Path + " Activating handler 'Placeholder'")
}

func InitRouter () *mux.Router {
// List of all routes for API
TheRoutes := []Route{
    Route {
        "Index",
        "Get",
        "/",
        Placeholder,
    },
}
router := mux.NewRouter()
fmt.Println("Loading Routes")
for _, r := range TheRoutes {
    router.HandleFunc(r.Path, r.Handler).Methods(r.Method)
}
return router
}

