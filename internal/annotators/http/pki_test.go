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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
)

func TestHttpPkiAnnotator_Do(t *testing.T) {
	b, err := ioutil.ReadFile("../../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req := httptest.NewRequest("POST", "/foo?param=value&foo=bar&baz=batman", nil)

	req.Header.Set("Host", "example.com")
	req.Header.Set("Date", "Tue, 20 Apr 2021 02:07:55 GMT")
	//req.Header.Set("Date", time.Now().String())
	// req.Header.Set("X-Example", "Example header with some whitespace.")
	// req.Header.Set("X-Empty-Header", "")
	// req.Header.Set("Cache-Control", "max-age=60")
	// req.Header.Add("Cache-Control", "must-revalidate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "18")

	//req.Header.Set("Signature-Input", "sig1=(\"@method\" \"@path\" \"@authority\" \"cache-control\" \"x-empty-header\" \"x-example\");created=1618884475 ;keyid=\"public\"; alg=\"ecdsa-p256-sha256\"")
	req.Header.Set("Signature-Input", "sig-b26=(\"date\" \"@method\" \"@path\" \"@authority\" \"content-type\" \"content-length\");created=1618884473;keyid=\"test-key-ed25519\"")

	req.Header.Set("Signature", "sig-b26=:wqcAqbmYJ2ji2glfAMaRy4gruYYnx2nEFN2HN6jrnDnQCK1u02Gb04v9EDgwUPiu4A0w6vuQv5lIp5WPpBKRCw==:")

	tpm := NewHttpPkiAnnotator(cfg)
	ctx := context.WithValue(req.Context(), "testData", req)

	anno, _ := tpm.Do(ctx, nil)
	t.Log(anno)
}
