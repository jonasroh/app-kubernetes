package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bem-vindo ao meu aplicativo de blog!")
	})
	http.ListenAndServe(":8080", nil)
}
