package logic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/google/uuid"
	"github.com/vatsal278/go-redis-cache"
	"github.com/vatsal278/htmltopdfsvc/internal/codes"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"github.com/vatsal278/htmltopdfsvc/internal/model"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/htmltopdfsvc/internal/logic HtmltopdfsvcLogicIer

type HtmltopdfsvcLogicIer interface {
	Ping(*model.PingRequest) *respModel.Response
	HtmlToPdf(w io.Writer, req *model.GenerateReq) *respModel.Response
	Upload(file multipart.File) *respModel.Response
	Replace(id string, file multipart.File) *respModel.Response
}

type htmltopdfsvcLogic struct {
	cacher redis.Cacher
}

func NewHtmltopdfsvcLogic(container *config.AppContainer) HtmltopdfsvcLogicIer {
	return &htmltopdfsvcLogic{
		cacher: container.Cacher,
	}
}

func (l htmltopdfsvcLogic) Ping(req *model.PingRequest) *respModel.Response {
	// add your business logic here
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "Pong",
		Data:    req.Data,
	}
}

func (l htmltopdfsvcLogic) HtmlToPdf(w io.Writer, req *model.GenerateReq) *respModel.Response {
	// add your business logic here
	var z map[string]interface{}
	b, err := l.cacher.Get(req.Id)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFetchingFile),
			Data:    nil,
		}
	}
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&z)
	if err != nil {
		log.Error("error unmarshaling JSON:" + err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileParseFail),
			Data:    nil,
		}
	}
	for i, p := range z["Pages"].([]interface{}) {
		page := p.(map[string]interface{})
		buf, err := base64.StdEncoding.DecodeString(page["Base64PageData"].(string))
		if err != nil {
			log.Error("error decoding base 64 input on page" + fmt.Sprint(i) + err.Error())
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrReadFileFail),
				Data:    nil,
			}
		}
		t, err := template.New(req.Id).Parse(string(buf))
		if err != nil {
			log.Error(err)
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrFileParseFail),
				Data:    nil,
			}
		}
		buffer := bytes.NewBuffer(nil)
		err = t.Execute(buffer, req.Values)
		if err != nil {
			log.Error(err)
			return &respModel.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrFileStoreFail),
				Data:    nil,
			}
		}
		page["Base64PageData"] = base64.StdEncoding.EncodeToString(buffer.Bytes())
	}
	buff := bytes.NewBuffer(nil)
	err = json.NewEncoder(buff).Encode(z)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewBuffer(buff.Bytes()))
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	pdfgFromJSON.SetOutput(w)
	err = pdfgFromJSON.Create()

	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data:    nil,
	}
}
func (l htmltopdfsvcLogic) Upload(file multipart.File) *respModel.Response {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrReadFileFail),
			Data:    nil,
		}
	}
	pdfg := wkhtmltopdf.NewPDFPreparer()
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(fileBytes)))
	jb, err := pdfg.ToJSON()
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileConversionFail),
			Data:    nil,
		}
	}
	u := uuid.NewString()
	err = l.cacher.Set(u, jb, 0)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data: map[string]interface{}{
			"id": u,
		},
	}

}

func (l htmltopdfsvcLogic) Replace(id string, file multipart.File) *respModel.Response {
	_, err := l.cacher.Get(id)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrKeyNotFound),
			Data:    nil,
		}
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrReadFileFail),
			Data:    nil,
		}
	}
	pdfg := wkhtmltopdf.NewPDFPreparer()
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(fileBytes)))
	jb, err := pdfg.ToJSON()
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileConversionFail),
			Data:    nil,
		}
	}

	err = l.cacher.Set(id, jb, 0)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFileStoreFail),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}
