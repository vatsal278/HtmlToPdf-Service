package handler

import (
	"bytes"
	"encoding/json"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvc(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *http.Request
		validateFunc func(*httptest.ResponseRecorder)
	}{
		{
			name:        "Success:: Upload",
			requestBody: "1",
			setupFunc: func() *http.Request {
				b := new(bytes.Buffer)
				y := multipart.NewWriter(b)
				part, _ := y.CreateFormFile("file", "./Failure.html")
				part.Write([]byte(`sample`))
				y.Close()
				r := httptest.NewRequest(http.MethodPost, "/v1/register", b)
				r.Header.Set("Content-Type", y.FormDataContentType())
				return r
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
		{
			name:        "Failure:: Upload:: no key found",
			requestBody: "1",
			setupFunc: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/v1/register", nil)
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setupFunc()
			w := httptest.NewRecorder()
			x := cacher.Upload
			x(w, r)
			tt.validateFunc(w)
		})
	}
}

func TestConvertToPdf(t *testing.T) {
	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvc(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() *http.Request
		validateFunc func(*httptest.ResponseRecorder)
	}{
		{
			name:        "Success:: ConvertToPdf",
			requestBody: "1",
			setupFunc: func() *http.Request {
				var temp struct {
					Values struct {
						Name  string `json:"name"`
						Marks int    `json:"marks"`
						ID    string `json:"id"`
					} `json:"values"`
				}
				temp.Values.Name = "vatsal"
				temp.Values.Marks = 90
				b, err := json.Marshal(temp)
				if err != nil {

				}
				r := httptest.NewRequest(http.MethodPost, "/v1/generate/eeab68db-acf0-46af-a6de-a08f909014e9", bytes.NewBuffer(b))
				return r
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
		{
			name:        "Failure:: ConvertToPdf",
			requestBody: "1",
			setupFunc: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/v1/register", nil)
			},
			validateFunc: func(x *httptest.ResponseRecorder) {
				if x.Code != http.StatusBadRequest {
					t.Errorf("want %v got %v", http.StatusBadRequest, x.Code)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setupFunc()
			w := httptest.NewRecorder()
			x := cacher.ConvertToPdf
			x(w, r)
			tt.validateFunc(w)
		})
	}
}
