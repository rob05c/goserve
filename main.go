package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"flag"
	"strconv"
)

var port int
func init() {
	const (
		defaultPort = 80
		portUsage = "the port on which to serve the website"
	)
	flag.IntVar(&port, "port", defaultPort, portUsage)
	flag.IntVar(&port, "p", defaultPort, portUsage + " (shorthand)")
	flag.Parse()
}

func main() {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	}

	http.HandleFunc("/", indexHandler)
	err = http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
