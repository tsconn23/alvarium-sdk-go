package utils

import (
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
	return true
}
