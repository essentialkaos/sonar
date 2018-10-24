package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"

	"pkg.re/essentialkaos/ek.v9/log"

	"github.com/valyala/fasthttp"

	"github.com/essentialkaos/sonar/slack"
	"github.com/essentialkaos/sonar/svg"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Points colors
const (
	COLOR_ONLINE = "#6DC185"
	COLOR_DND    = "#E7505A"
)

// Query strings
const (
	QUERY_MAIL  = "mail"
	QUERY_TOKEN = "token"
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

	statusHandler(ctx)
}

// requestRecover recover panic in request
func requestRecover(ctx *fasthttp.RequestCtx) {
	r := recover()

	if r != nil {
		log.Error("Recovered internal error in HTTP request handler: %v", r)
		ctx.SetStatusCode(501)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// statusHandler is status request handler
func statusHandler(ctx *fasthttp.RequestCtx) {
	query := ctx.QueryArgs()

	if !query.Has(QUERY_MAIL) || !query.Has(QUERY_TOKEN) {
		ctx.SetStatusCode(404)
		return
	}

	if bytes.Equal(token, query.Peek(QUERY_MAIL)) {
		ctx.SetStatusCode(404)
		return
	}

	ctx.Response.Header.Set("Content-Type", "image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "0")

	ctx.WriteString(getStatusBadge(string(query.Peek(QUERY_MAIL))))

	ctx.SetStatusCode(200)
}

// getStatusBadge return status badge
func getStatusBadge(mail string) string {
	if !enabled {
		return svg.GetBullet("")
	}

	// Bots always online
	if bots[mail] {
		return svg.GetBullet(COLOR_ONLINE)
	}

	switch slack.GetStatus(mail) {
	case slack.STATUS_OFFLINE:
		return svg.GetCircle()
	case slack.STATUS_ONLINE:
		return svg.GetBullet(COLOR_ONLINE)
	case slack.STATUS_DND:
		return svg.GetBullet(COLOR_DND)
	case slack.STATUS_VACATION:
		return svg.GetAirplane()
	case slack.STATUS_ONCALL:
		return svg.GetPhone()
	default:
		return svg.GetBullet("")
	}
}
