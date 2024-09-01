package main

import (
	"embed"
	"github.com/gorilla/handlers"
	"github.com/shiyanhui/dht/sample/web/routes"
	"log"
	"net/http"
)

func main() {
	startWebServer("2046")
}

//go:embed static
var static embed.FS

func startWebServer(port string) {
	r := routes.NewRouter()
	http.Handle("/", r)
	r.Handle("/static/", http.FileServer(http.FS(static)))
	//r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./sample/web/static/"))))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)) // Goroutine will block here

	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}
