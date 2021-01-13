package main

import (
	"net/http"

	"github.com/GreatLaboratory/go-web-example/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHTTPHandler())
}
