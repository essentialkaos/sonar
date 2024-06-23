module github.com/essentialkaos/sonar

go 1.19

replace github.com/slack-go/slack => ./pkgs/slack

require (
	github.com/essentialkaos/ek/v12 v12.127.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/slack-go/slack v0.13.0
	github.com/valyala/fasthttp v1.55.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/essentialkaos/depsy v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)
