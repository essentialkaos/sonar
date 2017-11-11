package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"pkg.re/essentialkaos/ek.v9/log"

	"github.com/valyala/fasthttp"

	"github.com/essentialkaos/sonar/slack"
	"github.com/essentialkaos/sonar/svg"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	COLOR_ERROR   = "#FFFFFF"
	COLOR_ONLINE  = "#6DC193"
	COLOR_DND     = "#E6D362"
	COLOR_OFFLINE = "#CCCCCC"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// startHTTPServer start HTTP server
func startHTTPServer(ip, port string) error {
	addr := ip + ":" + port

	log.Info("HTTP server is started on %s", addr)

	server := fasthttp.Server{
		Handler: fastHTTPHandler,
		Name:    APP + "/" + VER,
	}

	return server.ListenAndServe(addr)
}

// fastHTTPHandler handler for fast http requests
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	defer requestRecover(ctx)

	path := string(ctx.Path())

	if path != "/status.svg" {
		ctx.SetStatusCode(404)
		return
	}

	query := ctx.QueryArgs()

	if !query.Has("user") {
		ctx.SetStatusCode(404)
		return
	}

	writeBasicInfo(ctx)

	ctx.WriteString(getStatusBadge(string(query.Peek("user"))))
	ctx.SetStatusCode(200)
}

// requestRecover recover panic in request
func requestRecover(ctx *fasthttp.RequestCtx) {
	r := recover()

	if r != nil {
		log.Error("Recovered internal error in HTTP request handler: %v", r)
		ctx.SetStatusCode(501)
	}
}

// writeBasicInfo add basic info to response
func writeBasicInfo(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "0")
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getStatusBadge return status badge
func getStatusBadge(user string) string {
	if !enabled {
		return svg.GetPoint(COLOR_ERROR)
	}

	switch slack.GetStatus(user) {
	case slack.STATUS_OFFLINE:
		return svg.GetPoint(COLOR_OFFLINE)
	case slack.STATUS_ONLINE:
		return svg.GetPoint(COLOR_ONLINE)
	case slack.STATUS_DND:
		return svg.GetPoint(COLOR_DND)
	case slack.STATUS_VACATION:
		return svg.GetAirplane()
	default:
		return svg.GetPoint(COLOR_ERROR)
	}
}
