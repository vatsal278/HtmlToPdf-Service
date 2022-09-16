package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"
	"github.com/gorilla/mux"
	"github.com/vatsal278/htmltopdfsvc/internal/codes"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"github.com/vatsal278/htmltopdfsvc/internal/logic"
	"github.com/vatsal278/htmltopdfsvc/internal/model"
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
	ReplaceHtml(w http.ResponseWriter, r *http.Request)
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
	vars := mux.Vars(r)
	//we take id as a parameter from url path
	id, ok := vars["id"]
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNotfound), nil)
		return
	}
	var data model.GenerateReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrReadFileFail), nil)
		log.Error(err.Error())
		return
	}
	data.Id = id
	resp := svc.logic.HtmlToPdf(w, &data)
	w.Header().Set("Content-Disposition", "attachment; filename="+data.Id+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}

func (svc htmltopdfsvc) ReplaceHtml(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//we take id as a parameter from url path
	id, ok := vars["id"]
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrIdNotfound), nil)
		return
	}
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
	resp := svc.logic.Replace(id, file)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
