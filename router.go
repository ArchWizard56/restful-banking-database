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
    DualLogger(fmt.Sprintf("%s request from %s to %s Activating handler 'Placeholder'", r.Method,r.RemoteAddr,r.URL.Path))
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
DualLogger("Loading Routes")
for _, r := range TheRoutes {
    router.HandleFunc(r.Path, r.Handler).Methods(r.Method)
}
return router
}

