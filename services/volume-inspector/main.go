package main

import (
	"github.com/eudat-gef/gef/services/volume-inspector/api"
	"net/http"
)

func main() {
	http.ListenAndServe(":8181", api.Handlers())
}