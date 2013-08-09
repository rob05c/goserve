package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to " + r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}
