package main

import (
	"fmt"
	"net/http"

	"github.com/Ola-Daniel/qrcodebakery/internal/request"
	"github.com/Ola-Daniel/qrcodebakery/internal/response"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}


func (app *application) privacypolicy(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/privacypolicy.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) disclaimer(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/disclaimer.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}


func (app *application) tos(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/termsofservice.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

}


func (app *application) signup(w http.ResponseWriter, r *http.Request) {

}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {

    type response struct {
		DataString string `form:"dataString"`
	}
    var form response

	err := request.DecodePostForm(r, &form)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

     
	qrc, err := qrcode.New(form.DataString)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err) 
		return
	}
	
	wr, err := standard.New("./files/generated/test-qrcode.jpeg") 
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}
	
	// save file    
	if err = qrc.Save(wr); err != nil {
		fmt.Printf("could not save image: %v", err)
	}



}




func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
