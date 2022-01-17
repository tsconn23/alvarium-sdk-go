package annotators

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

	b, err := ioutil.ReadFile("../../test/res/config.json")
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

	ctx := context.WithValue(req.Context(), "Request", req)
	anno, _ := httppki.Do(ctx, nil)
	t.Log(anno)
}
