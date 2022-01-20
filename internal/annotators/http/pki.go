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
	"net/http"
	"os"

	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

// HttpPkiAnnotator is used to validate whether the signature on a given piece of data is valid
type HttpPkiAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewHttpPkiAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := HttpPkiAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationPKI
	a.sign = cfg.Signature
	return &a
}

func (a *HttpPkiAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {

	key := annotators.DeriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	//Call parser on request
	req := ctx.Value("Request")
	requestParser(req.(*http.Request))
	annotation := contracts.NewAnnotation(string(key), a.hash, hostname, a.kind, true)
	signed, err := annotators.SignAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(signed)
	return annotation, nil
}
