package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"

	"pkg.re/essentialkaos/ek.v12/log"

	"github.com/valyala/fasthttp"

	"github.com/essentialkaos/sonar/slack"
	"github.com/essentialkaos/sonar/svg"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Bullets colors
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

// startHTTPServer starts HTTP server
func startHTTPServer(ip, port string) error {
	addr := ip + ":" + port

	log.Info("HTTP server is started on %s", addr)

	server := fasthttp.Server{
		Handler: fastHTTPHandler,
		Name:    APP + "/" + VER,
	}

	return server.ListenAndServe(addr)
}

// fastHTTPHandler is handler for all requests
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	defer requestRecover(ctx)
	statusHandler(ctx)
}

// requestRecover recover panic in request
func requestRecover(ctx *fasthttp.RequestCtx) {
	r := recover()

	if r != nil {
		log.Error("Recovered internal error in HTTP request handler: %v", r)
		configureResponseHeaders(ctx)
		ctx.Write(svg.Empty)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// statusHandler is status request handler
func statusHandler(ctx *fasthttp.RequestCtx) {
	configureResponseHeaders(ctx)

	path := string(ctx.Path())

	if path != "/status.svg" {
		ctx.Write(svg.Empty)
		return
	}

	query := ctx.QueryArgs()

	if !query.Has(QUERY_MAIL) || !query.Has(QUERY_TOKEN) {
		ctx.Write(svg.Empty)
		return
	}

	if !bytes.Equal(token, query.Peek(QUERY_TOKEN)) {
		ctx.Write(svg.Empty)
		return
	}

	ctx.Write(getStatusBadge(query.Peek(QUERY_MAIL)))
}

// configureResponseHeaders configures response headers and status code
func configureResponseHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "0")
	ctx.SetStatusCode(200)
}

// getStatusBadge returns status badge
func getStatusBadge(userMail []byte) []byte {
	if !enabled {
		return svg.Offline
	}

	mail := string(userMail)

	// Bots always online
	if bots[mail] {
		return svg.GetBullet(COLOR_ONLINE)
	}

	switch slack.GetStatus(mail) {
	case slack.STATUS_OFFLINE:
		return svg.Offline
	case slack.STATUS_ONLINE:
		return svg.GetBullet(COLOR_ONLINE)
	case slack.STATUS_DND:
		return svg.GetBullet(COLOR_DND)
	case slack.STATUS_DND_OFFLINE:
		return svg.DND
	case slack.STATUS_VACATION:
		return svg.Vacation
	case slack.STATUS_ONCALL:
		return svg.OnCall
	case slack.STATUS_DISABLED:
		return svg.Empty
	default:
		return svg.Empty
	}
}
