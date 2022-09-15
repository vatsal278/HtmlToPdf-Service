package logic

import (
	"github.com/PereRohit/util/model"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
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
		//{
		//	name:        "Success:: Delete :: No key available",
		//	requestBody: "1",
		//	setupFunc: func() {
		//
		//	},
		//	validateFunc: func(x string) {
		//		if x != "del 1: 0" {
		//			t.Errorf("want %v got %v", "del 1: 0", x)
		//		}
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()

			x := cacher.Upload(data)
			tt.validateFunc(x)
		})
	}

}
