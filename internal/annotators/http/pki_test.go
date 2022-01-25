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
	"time"

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
	httppki := NewHttpPkiAnnotator(cfg)

	req := httptest.NewRequest("GET", "/path?param=value&foo=bar&baz=batman", nil)

	req.Header.Set("Host", "example.org")
	req.Header.Set("Date", time.Now().String())
	req.Header.Set("X-Example", "Example header with some whitespace.")
	req.Header.Set("X-Empty-Header", "")
	req.Header.Set("Cache-Control", "max-age=60")
	req.Header.Add("Cache-Control", "must-revalidate")

	req.Header.Set("Signature-Input", "sig1=(\"@method\" \"@path\" \"@authority\" \"cache-control\" \"x-empty-header\" \"x-example\");created=1618884475 ;keyid=\"test-key-rsa-pss\"; alg=\"my-algorithm\"")
	req.Header.Set("Signature", "sig1=:P0wLUszWQjoi54udOtydf9IWTfNhy+r53jGFj9XZuP4uKwxyJo1RSHi+oEF1FuX6O29d+lbxwwBao1BAgadijW+7O/PyezlTnqAOVPWx9GlyntiCiHzC87qmSQjvu1CFyFuWSjdGa3qLYYlNm7pVaJFalQiKWnUaqfT4LyttaXyoyZW84jS8gyarxAiWI97mPXU+OVM64+HVBHmnEsS+lTeIsEQo36T3NFf2CujWARPQg53r58RmpZ+J9eKR2CD6IJQvacn5A4Ix5BUAVGqlyp8JYm+S/CWJi31PNUjRRCusCVRj05NrxABNFv3r5S9IXf2fYJK+eyW4AiGVMvMcOg==:")

	ctx := context.WithValue(req.Context(), "Request", req)
	anno, _ := httppki.Do(ctx, nil)
	t.Log(anno)
}
