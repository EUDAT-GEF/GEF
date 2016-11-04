package main

import (
	"github.com/eudat-gef/gef/services/volume-inspector/api"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(":8181", api.Handlers())
}