package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func APITestCall(engine *gin.Engine, method, url string, body any, response any, apiToken string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, xerrors.Errorf("ERROR Marshal ", err)
	}
	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+apiToken)
	w.Header().Set("Accept", "application/json")
	req, _ := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	engine.ServeHTTP(w, req)

	logrus.Debug(req)
	if w.Code != 200 {
		return w.Result(), xerrors.Errorf("ERROR Code ", w.Result().Status)
	}
	if w.Body == nil {
		return w.Result(), xerrors.Errorf("ERROR Body is nil")
	}
	err = json.Unmarshal(w.Body.Bytes(), response)
	if err != nil {
		return w.Result(), xerrors.Errorf("ERROR Unmarshal ", err)
	}
	//logrus.Debugf("Response: %s", w.Body)
	return w.Result(), nil
}

func TestAPI(t *testing.T) {
	Init("../../conf/basic_strategy_dev_config.json", "../../conf/business_dev_config.json")
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestHealthz",
			func(t *testing.T) {
				var rssponse map[string]any
				_, err := APITestCall(Engine, "GET", "/api/healthz", nil, &rssponse, "")
				if err != nil {
					t.Error(err)
					return
				}
				t.Logf("Response: %v", rssponse)
			},
		},
		{
			"TestAuth",
			func(t *testing.T) {
				var request map[string]any
				request = make(map[string]any)
				request["apiKey"] = "string"
				_, err := APITestCall(Engine, "POST", "/api/auth", request, nil, "")
				if err != nil {
					t.Error(err)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}

}
