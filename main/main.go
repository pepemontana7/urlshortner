package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pepemontana7/urlshortner"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	pathsFilePtr := flag.String("paths-file", "paths.yaml", "Yaml or json paths file")
	flag.Parse()

	f, err := os.Open(*pathsFilePtr)
	check(err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)
	ft := strings.Split(*pathsFilePtr, ".")[1]

	var handler http.HandlerFunc
	if ft == "json" {
		fmt.Println("json handler")
		handler, err = urlshortner.JSONHandler(data, mapHandler)
	} else {
		fmt.Println("yaml handler")
		handler, err = urlshortner.YAMLHandler(data, mapHandler)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
