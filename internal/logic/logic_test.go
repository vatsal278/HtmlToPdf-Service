package logic

import (
	"github.com/PereRohit/util/model"
	"github.com/vatsal278/htmltopdfsvc/internal/codes"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	modelV "github.com/vatsal278/htmltopdfsvc/internal/model"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvcLogic(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() multipart.File
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Update",
			requestBody: "1",
			setupFunc: func() multipart.File {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body)
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()

			x := cacher.Upload(data)
			tt.validateFunc(x)
		})
	}

}

func TestHtmlToPdf(t *testing.T) {

	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvcLogic(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() modelV.GenerateReq
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Htmltopdf",
			requestBody: "1",
			setupFunc: func() modelV.GenerateReq {
				//body, _ := ioutil.TempFile(".", "example")
				var data modelV.GenerateReq
				data.Id = "ee5371c2-7200-45a1-b543-e4c5bd4c48ed"
				return data
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()
			var w io.Writer
			x := cacher.HtmlToPdf(w, &data)
			tt.validateFunc(x)
		})
	}

}

func TestReplace(t *testing.T) {
	os.Setenv("Address", "0.0.0.0")
	appContainer := config.GetAppContainer()
	cacher := NewHtmltopdfsvcLogic(appContainer)
	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (multipart.File, string)
		validateFunc func(*model.Response)
	}{
		{
			name:        "Success:: Replace",
			requestBody: "1",
			setupFunc: func() (multipart.File, string) {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body), "ee5371c2-7200-45a1-b543-e4c5bd4c48ed"
			},
			validateFunc: func(x *model.Response) {
				if x.Status != http.StatusCreated {
					t.Errorf("want %v got %v", http.StatusCreated, x)
				}
			},
		},
		{
			name:        "Failure:: Replace:: no key found",
			requestBody: "1",
			setupFunc: func() (multipart.File, string) {
				body, _ := ioutil.TempFile(".", "example")
				return multipart.File(body), ""
			},
			validateFunc: func(x *model.Response) {
				var temp *model.Response
				if !reflect.DeepEqual(x, &model.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrKeyNotFound),
					Data:    nil,
				}) {
					t.Errorf("want %v got %v", temp, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, id := tt.setupFunc()

			x := cacher.Replace(id, data)
			tt.validateFunc(x)
		})
	}
}
