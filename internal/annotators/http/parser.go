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
	"strings"
)

type parseResult struct {
	Seed      string
	Signature string
	Keyid     string
	Algorithm string
}

func RemoveExtraSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func parseRequest(r *http.Request) (parseResult, error) {
	//Signature Inputs extraction
	signatureInput := r.Header.Get("Signature-Input")
	signature := r.Header.Get("Signature")

	signatureInputList := strings.SplitN(signatureInput, ";", 2)

	signatureInputHeader := strings.Fields(signatureInputList[0])
	signatureInputTail := signatureInputList[1]

	var keyid, algorithm string

	signatureInputParsedTail := strings.Split(signatureInputTail, ";")
	for _, s := range signatureInputParsedTail {

		if strings.Contains(s, "alg") {
			raw := strings.Split(s, "=")[1]
			algorithm = strings.Trim(raw, "\"")
		}

		if strings.Contains(s, "keyid") {
			raw := strings.Split(s, "=")[1]
			keyid = strings.Trim(raw, "\"")
		}
	}

	signatureInputFields := make(map[string][]string)

	parsedSignatureInput := ""
	var s parseResult

	for _, field := range signatureInputHeader {
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
				rawQueryParams := strings.Split(r.URL.RawQuery, "&")
				var queryParams []string
				for _, rawQueryParam := range rawQueryParams {
					if rawQueryParam != "" {
						parameter := strings.Split(rawQueryParam, "=")
						name := parameter[0]
						value := parameter[1]
						b := new(bytes.Buffer)
						fmt.Fprintf(b, ";name=\"%s\": %s", name, value)
						queryParams = append(queryParams, b.String())
					}
				}
				signatureInputFields[key] = queryParams
			default:
				return s, fmt.Errorf("Unhandled Specialty Component %s", key)
			}
		} else {
			fieldValues := r.Header.Values(key)

			if len(fieldValues) == 0 {
				return s, fmt.Errorf("Unhandled Specialty Component %s", key)
			} else if len(fieldValues) == 1 {
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
			parsedSignatureInput += ("\"" + key + "\" " + keyValues[0] + "\n")
		} else {
			for _, v := range keyValues {
				parsedSignatureInput += ("\"" + key + "\"" + v + "\n")
			}
		}
	}

	parsedSignatureInput = fmt.Sprintf("%s;%s", parsedSignatureInput, signatureInputTail)
	s = parseResult{Seed: parsedSignatureInput, Signature: signature, Keyid: keyid, Algorithm: algorithm}

	return s, nil
}
