
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"errors"
)

var port int
var log string
var files string
func init() {
	const (
		defaultLog  = ""
		logUsage    = "the log file"
		defaultPort = 80
		portUsage   = "the port on which to serve the website"
		defaultFiles = "files.json"
		filesUsage = "the file which contains a json object of file names to serve and their content types"
	)
	flag.IntVar(&port, "port", defaultPort, portUsage)
	flag.IntVar(&port, "p", defaultPort, portUsage+" (shorthand)")
	flag.StringVar(&log, "log", defaultLog, logUsage)
	flag.StringVar(&log, "l", defaultLog, logUsage+" (shorthand)")
	flag.StringVar(&files, "files", defaultFiles, filesUsage)
	flag.StringVar(&files, "f", defaultFiles, filesUsage+" (shorthand)")
	flag.Parse()
}

type ServePath struct {
	Path        string
	Value       []byte
	ContentType string
}

func NewFileServePath(file string, contentType string) (*ServePath, error) {
	value, err := ioutil.ReadFile(file)
	return &ServePath{Path: file, ContentType: contentType, Value: value}, err
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

func parseFiles(j []byte) (map[string]string, error) {
	var raw interface{}
	err := json.Unmarshal(j, &raw)
	if err != nil {
		return nil, err
	}
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.New("file was not JSON object")
	}
	fileMap := make(map[string]string)
	for k, v := range rawMap {
		vs, ok := v.(string)
		if !ok {
			continue
		}
		fileMap[k] = vs
	}
	return fileMap, nil
}

func main() {
	logChan := NewLogFileWriter(log)

	jsonFile, err := ioutil.ReadFile(files)
	if err != nil {
		fmt.Println("Could not find list file '" + files + "'. Nothing to serve.")
		return
	}
	fileMap, err := parseFiles(jsonFile)
	if err != nil {
		fmt.Println("Error reading list file. Nothing to serve.")
		fmt.Println(err)
		return
	}
	paths := make(map[string]*ServePath, 0)
	for file, contentType := range fileMap {
		path, err := NewFileServePath(file, contentType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paths[file] = path
	}

	makeHandler := func(s []byte, contentType string) func(w http.ResponseWriter, r *http.Request) {
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
		http.HandleFunc("/"+path.Path, makeHandler(path.Value, path.ContentType))
	}
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
	}
	close(logChan)
}
