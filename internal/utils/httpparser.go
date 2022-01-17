package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

//Returns a boolean for the moment
func HTTPParser(r *http.Request) bool {

	//Signature Inputs extraction
	var rgx = regexp.MustCompile(`\((.*?)\)`)
	signature_input := r.Header.Get("Signature-Input")
	rs := rgx.FindStringSubmatch(signature_input)
	signature_input_list := strings.Split(rs[1], " ")
	fmt.Printf("Signature-Input list: %v\n", signature_input_list)
	speciality_components := make(map[string][]string)

	//@method:
	fmt.Printf(" --->  @method: %v\n", r.Method)
	speciality_components["@method"] = []string{r.Method}

	//@authority:
	fmt.Printf(" --->  @authority: %v\n", r.Host)
	speciality_components["@authority"] = []string{r.Host}

	//@scheme:
	protool := r.Proto
	scheme := strings.ToLower(strings.Split(protool, "/")[0])

	fmt.Printf(" --->  @scheme : %v\n", scheme)
	speciality_components["@scheme"] = []string{scheme}

	//@request-target
	fmt.Printf(" --->  @request-target: %v\n", r.RequestURI)
	speciality_components["@request-target"] = []string{r.RequestURI}

	//@path
	fmt.Printf(" --->  @path: %v\n", r.URL.Path)
	speciality_components["@path"] = []string{r.URL.Path}

	//@query
	var query string = "?"
	query += r.URL.RawQuery

	fmt.Printf(" --->  @query: %v\n", query)
	speciality_components["@query"] = []string{query}

	//@query-params
	query_params_raw_map := r.URL.Query()
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
