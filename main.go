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

type ServePath struct {
	Path string
	Value []byte
	ContentType string
}

func NewFileServePath(file string, contentType string) (*ServePath, error) {
	value, err := ioutil.ReadFile(file)
	return &ServePath {Path: file, ContentType: contentType, Value: value}, err
}

func main() {
	// only specific files are served. A user can't just request an arbitrary file.
	// @todo put this in a config file
	files := map[string] string {
		"index.html": "text/html", // the first entry will be served at the root
		"main.css": "text/css",
		"reset.css": "text/css",
	}
	paths := make([]*ServePath, 0)
	for file, contentType := range files {
		path, err := NewFileServePath(file, contentType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paths = append(paths, path)
	}

	makeHandler := func(s []byte, contentType string) (func(w http.ResponseWriter, r *http.Request)) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			fmt.Fprintf(w, "%s", s)
		}
	}

	if len(paths) == 0 {
		fmt.Println("No paths to serve")
		return
	}

	http.HandleFunc("/", makeHandler(paths[0].Value, paths[0].ContentType))
	for _, path := range paths {
		http.HandleFunc("/" + path.Path, makeHandler(path.Value, path.ContentType))
	}
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
