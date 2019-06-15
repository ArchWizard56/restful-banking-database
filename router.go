package main
import (
    "github.com/gorilla/mux"
)
type Route struct {
    Name string
    Method string
    Path string
    Handler string
}
func tmp() []Route{
TheRoutes := []Route{
    Route {
        "Index",
        "Get",
        "/",
        "test",
    },
}
return TheRoutes
}
