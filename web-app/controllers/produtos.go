package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pbnjay/memory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"main.go/models"
)

/*
Must e qm encapsula todos os templates
ParseGlob para especificar o caminho
*/

func MetricsHandler(w http.ResponseWriter, r *http.Request) {

	monitorandoMemoria()
	monitorandoCPU()

	// Coleta e exporta as métricas
	promhttp.Handler().ServeHTTP(w, r)
}

var temp = template.Must(template.ParseGlob("templates/*.html"))

func memoriaLivre() float64 { // func para pegar a memoria livre
	memoria_livre := memory.FreeMemory()
	return float64(memoria_livre)
}

func totalMemory() float64 { // func pegar a memória total
	memoria_total := memory.TotalMemory()
	return float64(memoria_total)
}

var (
	latencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_latency_seconds",
		Help:    "Histogram of the request latency in seconds.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})

	memoriaLivreBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memoria_livre_bytes",                  // nome da métrica
		Help: "Quantidade de memória livre em bytes", // descricão de métrica
	})

	memoriaLivreMegabytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memoria_livre_megabytes",                  // nome da métrica
		Help: "Quantidade de memória livre em megabytes", // descricão de métrica
	})

	totalMemoryBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_memoria_bytes",                  // nome da métrica
		Help: "Quantidade total de memória em bytes", // descricão de métrica
	})

	totalMemoryGigaBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_memoria_gigabytes",                  // nome da métrica
		Help: "Quantidade total de memória em gigabytes", // descricão de métrica
	})

	cpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "uso_cpu",
		Help: "Uso da CPU em porcentagem",
	},
		[]string{"cpu"},
	)
)

func init() { // funcão para registrar as métricas
	prometheus.MustRegister(latencyHistogram)
	prometheus.MustRegister(memoriaLivreBytesGauge)
	prometheus.MustRegister(memoriaLivreMegabytesGauge)
	prometheus.MustRegister(totalMemoryBytesGauge)
	prometheus.MustRegister(totalMemoryGigaBytesGauge)
	prometheus.MustRegister(cpuUsage)
}

func monitorandoMemoria() {
	memoriaLivreBytesGauge.Set(memoriaLivre())
	memoriaLivreMegabytesGauge.Set(memoriaLivre() / 1024 / 1024)
	totalMemoryBytesGauge.Set(totalMemory())
	totalMemoryGigaBytesGauge.Set(totalMemory() / 1024 / 1024 / 1024)
}

func monitorandoCPU() {
	go func() {
		for {
			usage, err := cpu.Percent(0, true)
			if err == nil {
				for i, u := range usage {
					cpuUsage.WithLabelValues(fmt.Sprintf("cpu%d", i)).Set(u)
				}
			}

			time.Sleep(15 * time.Second)
		}
	}()
}

func Index(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	todosOsProdutos := models.BuscaTodosOsProdutos()
	temp.ExecuteTemplate(w, "Index", todosOsProdutos)

	latency := time.Since(start).Seconds()
	latencyHistogram.Observe(latency)
}

func New(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nome := r.FormValue("nome")
		descricao := r.FormValue("descricao")
		preco := r.FormValue("preco")
		quantidade := r.FormValue("quantidade")

		precoConvertidoParaFloat, err := strconv.ParseFloat(preco, 64)
		if err != nil {
			log.Println("Erro na conversão do preço:", err)
		}
		quantidadeConvertidaParaInt, err := strconv.Atoi(quantidade)
		if err != nil {
			log.Println("Erro na conversão da quantidade:", err)
		}

		models.CriarNovoProduto(nome, descricao, precoConvertidoParaFloat, quantidadeConvertidaParaInt)

	}
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idDoProduto := r.URL.Query().Get("id")
	models.DeletaProduto(idDoProduto)
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idDoProduto := r.URL.Query().Get("id")
	produto := models.ExibeProduto(idDoProduto)
	temp.ExecuteTemplate(w, "Edit", produto)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nome := r.FormValue("nome")
		descricao := r.FormValue("descricao")
		preco := r.FormValue("preco")
		quantidade := r.FormValue("quantidade")

		idConvertidoParaInt, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Erro na conversão do ID para Int:", err)
		}

		precoConvertidoParaFloat, err := strconv.ParseFloat(preco, 64)
		if err != nil {
			log.Println("Erro na conversão do preço para float64:", err)
		}

		quantidadeConvertidaParaInt, err := strconv.Atoi(quantidade)
		if err != nil {
			log.Println("Erro na conversão da quantidade para Int:", err)
		}

		models.AtualizaProduto(idConvertidoParaInt, nome, descricao, precoConvertidoParaFloat, quantidadeConvertidaParaInt)
	}
	http.Redirect(w, r, "/", 301)
}
