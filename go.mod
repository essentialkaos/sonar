module github.com/essentialkaos/sonar

go 1.19

replace github.com/slack-go/slack => ../slack

require (
	github.com/essentialkaos/ek/v12 v12.62.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/slack-go/slack v0.12.1
	github.com/valyala/fasthttp v1.45.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
