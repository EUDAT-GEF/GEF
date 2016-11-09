package main

import (
	"github.com/EUDAT-GEF/GEF/services/volume-inspector/api"
	"net/http"
)

func main() {
	http.ListenAndServe(":8282", api.Handlers())
}