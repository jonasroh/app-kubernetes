package routes

import (
	"net/http"

	// biblioteca criar o servidor web
	"main.go/controllers"
)

func CarregaRotas() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/new", controllers.New)
	http.HandleFunc("/insert", controllers.Insert)
	http.HandleFunc("/delete", controllers.Delete)
	http.HandleFunc("/edit", controllers.Edit)
	http.HandleFunc("/update", controllers.Update)

	http.Handle("/metrics", http.HandlerFunc(controllers.MetricsHandler)) // criando o servidor web para expor as métricas

}
