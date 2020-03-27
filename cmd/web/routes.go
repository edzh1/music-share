package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	// standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/link", http.HandlerFunc(app.handleLink))

	// return standardMiddleware.Then(mux)
	return mux
}
