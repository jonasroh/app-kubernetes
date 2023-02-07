package main

import (
	"net/http"
	"text/template"
)

type Produto struct {
	Nome       string
	Descricao  string
	Preco      float64
	Quantidade int
}

/*
Must e qm encapsula todos os templates
ParseGlob para especificar o caminho
*/
var temp = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	produtos := []Produto{
		{Nome: "Camiseta", Descricao: "Manga longa", Preco: 49, Quantidade: 5},
		{"Camiseta", "Oversized", 79, 3},
		{"Calça", "Jeans", 109, 7},
	}
	temp.ExecuteTemplate(w, "Index", produtos)
}
