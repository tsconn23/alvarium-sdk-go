/*******************************************************************************
 * Copyright 2021 Dell Inc.
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
package annotators

import (
	"context"
	"os"

	"github.com/project-alvarium/alvarium-sdk-go/internal/utils"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

// PkiAnnotator is used to validate whether the signature on a given piece of data is valid
type HTTPPkiAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewHttpPkiAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := HTTPPkiAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationPKI
	a.sign = cfg.Signature
	return &a
}

func (a *HTTPPkiAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {

	key := deriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	// var sig signable
	// err := json.Unmarshal(data, &sig)
	// if err != nil {
	// 	return contracts.Annotation{}, err
	// }

	//Call parser on request
	req := ctx.Value("Request")
	utils.HTTPParser(req)
	annotation := contracts.NewAnnotation(string(key), a.hash, hostname, a.kind, true)
	signed, err := signAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(signed)
	return annotation, nil
}
