package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"
	"github.com/vatsal278/htmltopdfsvc/internal/codes"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"github.com/vatsal278/htmltopdfsvc/internal/logic"
	"github.com/vatsal278/htmltopdfsvc/internal/model"
	"io/ioutil"
	"net/http"
)

const HtmltopdfsvcName = "htmltopdfsvc"

var (
	svc = &htmltopdfsvc{}
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/htmltopdfsvc/internal/handler HtmltopdfsvcHandler

type HtmltopdfsvcHandler interface {
	HealthChecker
	Ping(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	ConvertToPdf(w http.ResponseWriter, r *http.Request)
}

type htmltopdfsvc struct {
	logic logic.HtmltopdfsvcLogicIer
}

func init() {
	AddHealthChecker(svc)
}

func NewHtmltopdfsvc(c *config.AppContainer) HtmltopdfsvcHandler {
	svc = &htmltopdfsvc{
		logic: logic.NewHtmltopdfsvcLogic(c),
	}

	return svc
}

func (svc htmltopdfsvc) HealthCheck() (svcName string, msg string, stat bool) {
	set := false
	defer func() {
		svcName = HtmltopdfsvcName
		if !set {
			msg = ""
			stat = true
		}
	}()
	return
}

func (svc htmltopdfsvc) Ping(w http.ResponseWriter, r *http.Request) {
	req := &model.PingRequest{}

	suggestedCode, err := request.FromJson(r, req)
	if err != nil {
		response.ToJson(w, suggestedCode, fmt.Sprintf("FAILED: %s", err.Error()), nil)
		return
	}
	// call logic
	resp := svc.logic.Ping(req)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
	return
}

func (svc htmltopdfsvc) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10240) //File size to come from config
	if err != nil {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileSizeExceeded), nil)
		log.Error(err.Error())
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
		log.Error(err.Error())
		return
	}
	defer file.Close()
	resp := svc.logic.Upload(file)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}

func (svc htmltopdfsvc) ConvertToPdf(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("file")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrFileParseFail), nil)
		log.Error(err.Error())
	}
	var class model.Class
	json.Unmarshal(data, &class)
	resp := svc.logic.HtmlToPdf(v, class)
	w.Header().Set("Content-Disposition", "attachment; filename="+v+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
