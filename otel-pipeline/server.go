package main

import (
	"fmt"
	"net/http"
)

var visitCounter = 0

func updateVisitCounter(writer http.ResponseWriter, req *http.Request) {
	visitCounter++
	fmt.Fprintf(writer, "This page has been visited %v times.\n", visitCounter)
}

func getVisitCounter(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "%v", visitCounter)
}

func main() {
	http.HandleFunc("/visit", updateVisitCounter)
	http.HandleFunc("/getVisitCounter", getVisitCounter)
	http.ListenAndServe(":8090", nil)
}
