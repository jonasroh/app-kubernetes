package main

import (
	"net/http"

	"main.go/routes"
)

func main() {
	routes.CarregaRotas()
	http.ListenAndServe(":8000", nil)

}

/* modularizado de acordo com o mvc */
