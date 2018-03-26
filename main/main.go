package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pepemontana7/urlshortner"
	. "github.com/pepemontana7/urlshortner/dao"
	"gopkg.in/mgo.v2/bson"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	pathsFilePtr := flag.String("paths-file", "", "Yaml or json paths file")
	flag.Parse()
	var ft string
	var data []byte
	if *pathsFilePtr != "" {
		ft = strings.Split(*pathsFilePtr, ".")[1]
		f, err := os.Open(*pathsFilePtr)
		check(err)
		defer f.Close()

		data, err = ioutil.ReadAll(f)
		check(err)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mongoPaths := map[string]string{
		"/mongo-urlshort": "https://godoc.org/github.com/gophercises/urlshort",
		"/mongo-godoc":    "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	var handler http.HandlerFunc
	var err error
	if ft == "json" {
		fmt.Println("json handler")
		handler, err = urlshortner.JSONHandler(data, mapHandler)
	} else if ft == "yaml" {
		fmt.Println("yaml handler")
		handler, err = urlshortner.YAMLHandler(data, mapHandler)
	} else {
		//connect to db
		dao := PathsDAO{}
		dao.Server = "localhost"
		dao.Database = "pathsdb"
		dao.Connect()
		// create paths in db from pathsToURL
		for path, url := range mongoPaths {
			if _, err := dao.FindByPath(path); err != nil {
				fmt.Printf("Could not find %s so adding", path)
				var p Path
				p.ID = bson.NewObjectId()
				p.Path = path
				p.URL = url
				if err := dao.Insert(p); err != nil {
					fmt.Println("error creating, ", err)
					panic(err)
				}
			}

		}
		s, e := dao.FindAll()
		check(e)
		fmt.Println(s)
		handler, err = urlshortner.DBHandler(&dao, mapHandler)

	}
	check(err)
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
