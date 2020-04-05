package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", hub.handleWebSocket)
	http.HandleFunc("/", serveIndex)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("client/index.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
