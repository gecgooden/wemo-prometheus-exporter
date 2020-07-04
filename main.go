package main

import (
	"log"
	"net/http"

	arg "github.com/alexflint/go-arg"
	wemo "github.com/gecgooden/go.wemo"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config gets its content from env and passes it on to different packages
type Config struct {
	WemoHost string `arg:"env:HOST"`
	WebAddr  string `arg:"env:WEB_ADDR"`
	WebPath  string `arg:"env:WEB_PATH"`
}

func main() {
	log.Println("Starting wemo-prometheus-exporter")

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env present")
	}

	c := Config{
		WemoHost: "http://localhost",
		WebPath:  "/metrics",
		WebAddr:  ":8080",
	}

	arg.MustParse(&c)

	device := &wemo.Device{
		Host: c.WemoHost,
	}

	prometheus.MustRegister(NewWemoCollector(device))

	http.Handle(c.WebPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Wemo Exporter</title></head>
			<body>
			<h1>Wemo Exporter</h1>
			<p><a href="` + c.WebPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(c.WebAddr, nil))
}
