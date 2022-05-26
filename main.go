package main

import (
	"fmt"
	"log"
	"net/http"

	"blupine.co/gmail-cleaner/routers"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	router := routers.GmailRouter{}
	http.HandleFunc("/flush", router.FlushMessages)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
