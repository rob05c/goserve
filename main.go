package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"flag"
	"strconv"
	"os"
	"time"
)

var port int
var log string
func init() {
	const (
		defaultLog = ""
		logUsage = "the log file"
		defaultPort = 80
		portUsage = "the port on which to serve the website"
	)
	flag.IntVar(&port, "port", defaultPort, portUsage)
	flag.IntVar(&port, "p", defaultPort, portUsage + " (shorthand)")
	flag.StringVar(&log, "log", defaultLog, logUsage)
	flag.StringVar(&log, "l", defaultLog, logUsage + " (shorthand)")
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

func NewLogFileWriter(logFile string) chan string {
	if logFile == "" {
		return nil
	}
	logChan := make(chan string, 100)

	megabyte := int64(1048576)
	fileMax := megabyte * 100

	fi, err := os.Create(logFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	go func() {
		defer fi.Close()
		for {
			
			stat, err := fi.Stat()
			if err != nil && stat != nil && stat.Size() > fileMax {
				fi.Seek(0, 0)
				fi.Truncate(0)
				fi.WriteString("(truncated)\n")
			}
			msg, ok := <-logChan
			if !ok {
				break
			}
			_, err = fi.WriteString(msg)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}()
	return logChan
}

func main() {
	logChan := NewLogFileWriter(log)

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
			if logChan != nil {
				logChan <- time.Now().String() + " " + r.RemoteAddr + "\n"
			}
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
	close(logChan)
}
