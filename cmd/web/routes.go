package main

import (
	"net/http"

	"github.com/Ola-Daniel/qrcodebakery/assets"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()
	mux.NotFoundHandler = http.HandlerFunc(app.notFound)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.PathPrefix("/static/").Handler(fileServer)

	mux.HandleFunc("/", app.home).Methods("GET")

	protectedRoutes := mux.NewRoute().Subrouter()
	protectedRoutes.Use(app.requireBasicAuthentication)
	protectedRoutes.HandleFunc("/basic-auth-protected", app.protected).Methods("GET")

	return mux
}
