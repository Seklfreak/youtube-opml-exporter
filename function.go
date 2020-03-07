package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Seklfreak/youtube-opml-exporter/pkg"
)

func main() {
	http.HandleFunc("/", pkg.ExportHandler)
	http.HandleFunc("/login", pkg.LoginHandler)
	http.HandleFunc("/exchange", pkg.ExchangeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
