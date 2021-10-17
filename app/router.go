package main

import (
	"net/http"

	"github.com/corioders/terrain-story.api/app/routes/qr"
	"github.com/corioders/terrain-story.api/foundation"

	"github.com/corioders/gokit/web"
	"github.com/corioders/gokit/web/middleware"
)

func newRouter(app *foundation.Application) (http.Handler, error) {
	routerLogger := app.GetLogger().Child("Router")
	_ = routerLogger

	router := web.NewRouter(app.GetLogger(),
		middleware.Errors(app.GetLogger()),
		middleware.Compression(),
		middleware.Cors("*"),
	)

	qrController, err := qr.NewController(app.GetConfig().Qr)
	if err != nil {
		return nil, err
	}

	router.Handle(http.MethodGet, "/qr/:uuid", qrController.RedirectHandler)

	return router, nil
}
