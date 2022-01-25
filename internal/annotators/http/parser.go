/*******************************************************************************
 * Copyright 2022 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package http

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func RemoveExtraSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func requestParser(r *http.Request) string {

	//Signature Inputs extraction
	var rgx = regexp.MustCompile(`\((.*?)\)`)
	fmt.Printf("r: %v\n", r)
	signatureInput := r.Header.Get("Signature-Input")
	fmt.Printf("signatureInput: %v\n", signatureInput)

	rs := rgx.FindStringSubmatch(signatureInput)
	signatureInputList := strings.Split(rs[1], " ")
	fmt.Printf("Signature-Input list: %v\n", signatureInputList)

	signatureInputFields := make(map[string][]string)
	parsedSignatureInput := ""

	signatureInputParsedSection := strings.Split(signatureInput, ";")

	fmt.Printf("---> signatureInputSections: %v\n", signatureInputParsedSection)

	for _, s := range signatureInputParsedSection {

		if strings.Contains(s, "alg") {
			algorthim_raw := strings.Split(s, "=")[1]
			algorthim := strings.Trim(algorthim_raw, "\"")
			signatureInputFields["alg"] = []string{algorthim}
		}

		if strings.Contains(s, "key") {
			keyid_raw := strings.Split(s, "=")[1]
			keyid := strings.Trim(keyid_raw, "\"")
			signatureInputFields["keyid"] = []string{keyid}
		}

	}

	// Now we have the value of the keyid and algorithm in the signatureInputFields
	// the next line is for logging
	fmt.Printf("=============>> sig input ---------%v\n", signatureInputFields)

	for _, field := range signatureInputList {
		//remove double quotes from the field to access it directly in the header map
		key := field[1 : len(field)-1]
		if key[0:1] == "@" {
			switch specialtyComponent(key) {
			case method:
				signatureInputFields[key] = []string{r.Method}
			case authority:
				signatureInputFields[key] = []string{r.Host}
			case scheme:
				protool := r.Proto
				scheme := strings.ToLower(strings.Split(protool, "/")[0])
				signatureInputFields[key] = []string{scheme}
			case requestTarget:
				signatureInputFields[key] = []string{r.RequestURI}
			case path:
				signatureInputFields[key] = []string{r.URL.Path}
			case query:
				var query string = "?"
				query += r.URL.RawQuery
				signatureInputFields[key] = []string{query}
			case queryParams:
				queryParamsRawMap := r.URL.Query()
				var queryParams []string
				for key, value := range queryParamsRawMap {
					b := new(bytes.Buffer)
					fmt.Fprintf(b, ";name=\"%s\": %s", key, value[0])
					queryParams = append(queryParams, b.String())
				}

				signatureInputFields[key] = queryParams
			default:
				fmt.Println("Unhandled Specialty Component")
			}
		} else {
			fieldValues := r.Header.Values(key)

			if len(fieldValues) == 1 {
				value := RemoveExtraSpaces(r.Header.Get(key))
				signatureInputFields[key] = []string{value}

			} else {

				value := ""
				for i := 0; i < len(fieldValues); i++ {
					value += fieldValues[i]
					if i != (len(fieldValues) - 1) {
						value += ", "
					}
				}
				value = RemoveExtraSpaces(value)
				signatureInputFields[key] = []string{value}
			}
		}
		// Construct final output string
		keyValues := signatureInputFields[key]
		if len(keyValues) == 1 {
			parsedSignatureInput += ("\"" + key + "\": " + keyValues[0] + "\n")
		} else {
			for _, v := range keyValues {
				parsedSignatureInput += ("\"" + key + "\"" + v + "\n")
			}
		}
	}

	//remove signature name
	index := strings.Index(signatureInput, "=")
	signatureInput = signatureInput[index+1:]

	// check if new line needs to be added at the end
	parsedSignatureInput += ("\"@signature-params\": " + signatureInput + "\n")
	fmt.Printf("FINAL OUTPUT %v\n", parsedSignatureInput)

	return parsedSignatureInput
}
