package urlshortner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pepemontana7/urlshortner/dao"
	"gopkg.in/yaml.v2"
)

type paths struct {
	Path string
	Url  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Map Handler", r.RequestURI)
		u, ok := pathsToUrls[r.RequestURI]
		fmt.Println(u)
		if !ok {
			fmt.Println("Map Handler: no redirect for", r.RequestURI)
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, u, 301)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	parsedJson, err := parseJSON(j)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}
func DBHandler(d *dao.PathsDAO, fallback http.Handler) (http.HandlerFunc, error) {
	fmt.Println("DB Handler..")

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("DB Handler", r.RequestURI)
		path, err := d.FindByPath(r.RequestURI)
		if err != nil {
			fmt.Println("error infind by path")
			fallback.ServeHTTP(w, r)
			return
		}
		fmt.Println("db path ", path)
		http.Redirect(w, r, path.URL, 301)
	}, nil
}

func parseYAML(y []byte) ([]paths, error) {
	//fmt.Println("yml: data ", string(y))

	var b []paths

	err := yaml.Unmarshal(y, &b)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
		return nil, err
	}

	for _, i := range b {
		fmt.Println("yml: ", i, i.Url, i.Path)
	}
	return b, nil
}

func parseJSON(j []byte) ([]paths, error) {
	var b []paths
	fmt.Println(b)
	err := json.Unmarshal(j, &b)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
		return nil, err
	}
	for _, i := range b {
		fmt.Println("json: ", i, i.Url, i.Path)
	}
	return b, nil
}

func buildMap(py []paths) map[string]string {
	mpy := map[string]string{}
	for _, y := range py {
		mpy[y.Path] = y.Url
	}
	fmt.Println("mpy: ", mpy)
	return mpy
}
