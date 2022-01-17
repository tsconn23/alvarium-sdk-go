package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

//Returns a boolean for the moment
func HTTPParser(r interface{}) bool {

	speciality_components := make(map[string][]string)

	//@method:
	fmt.Printf(" --->  @method: %v\n", r.(*http.Request).Method)
	speciality_components["@method"] = []string{r.(*http.Request).Method}

	//@authority:
	fmt.Printf(" --->  @authority: %v\n", r.(*http.Request).Host)
	speciality_components["@authority"] = []string{r.(*http.Request).Host}

	//@scheme:
	protool := r.(*http.Request).Proto
	scheme := strings.ToLower(strings.Split(protool, "/")[0])

	fmt.Printf(" --->  @scheme : %v\n", scheme)
	speciality_components["@scheme"] = []string{scheme}

	//@request-target
	fmt.Printf(" --->  @request-target: %v\n", r.(*http.Request).RequestURI)
	speciality_components["@request-target"] = []string{r.(*http.Request).RequestURI}

	//@path
	fmt.Printf(" --->  @path: %v\n", r.(*http.Request).URL.Path)
	speciality_components["@path"] = []string{r.(*http.Request).URL.Path}

	//@query
	var query string = "?"
	query += r.(*http.Request).URL.RawQuery

	fmt.Printf(" --->  @query: %v\n", query)
	speciality_components["@query"] = []string{query}

	//@query-params
	query_params_raw_map := r.(*http.Request).URL.Query()
	var query_params []string
	for key, value := range query_params_raw_map {
		b := new(bytes.Buffer)
		fmt.Fprintf(b, " ;name=\"%s\":%s", key, value[0])
		query_params = append(query_params, b.String())
	}

	fmt.Printf(" --->  @query-params: %v\n", query_params)
	speciality_components["@query"] = query_params

	//printing entire map
	fmt.Printf(" --->  @speciality_components: %v\n", speciality_components)

	return true
}
