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
		"main.js": "text/javascript",
		"home.txt": "application/octet-stream",
		"projects.txt": "application/octet-stream",
		"resume.txt": "application/octet-stream",
		"contact.txt": "application/octet-stream",
		"about.txt": "application/octet-stream",
		"license.txt": "application/octet-stream",
		"favicon.ico": "image/x-icon",
	}
	paths := make(map[string]*ServePath, 0)
	for file, contentType := range files {
		path, err := NewFileServePath(file, contentType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paths[file] = path
	}

	makeHandler := func(s []byte, contentType string) (func(w http.ResponseWriter, r *http.Request)) {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("serving " + r.URL.Path)
			w.Header().Set("Content-Type", contentType)
			fmt.Fprintf(w, "%s", s)
		}
	}

	if len(paths) == 0 {
		fmt.Println("No paths to serve")
		return
	}

	root := paths["index.html"]
	http.HandleFunc("/", makeHandler(root.Value, root.ContentType))
	for _, path := range paths {
		http.HandleFunc("/" + path.Path, makeHandler(path.Value, path.ContentType))
	}
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
