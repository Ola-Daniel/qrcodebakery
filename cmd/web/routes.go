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
	//mux.PathPrefix("/files/").Handler(fileServer)  


	mux.HandleFunc("/", app.home).Methods("GET")

	mux.HandleFunc("/tos", app.tos).Methods("GET")

	mux.HandleFunc("/privacy-policy", app.privacypolicy).Methods("GET")

	mux.HandleFunc("/disclaimer", app.disclaimer).Methods("GET")

	mux.HandleFunc("/generate", app.generate).Methods("POST")

    mux.HandleFunc("/login", app.login).Methods("GET", "POST")

	mux.HandleFunc("/sign-up", app.signup).Methods("GET", "POST")



	protectedRoutes := mux.NewRoute().Subrouter()
	protectedRoutes.Use(app.requireBasicAuthentication)
	protectedRoutes.HandleFunc("/basic-auth-protected", app.protected).Methods("GET")

	return mux
}
