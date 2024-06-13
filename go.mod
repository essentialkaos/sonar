module github.com/essentialkaos/sonar

go 1.19

replace github.com/slack-go/slack v0.12.5 => github.com/essentialkaos/slack v0.0.0-20240328220008-c8788ecc08aa

require (
	github.com/essentialkaos/ek/v12 v12.126.1
	github.com/orcaman/concurrent-map v1.0.0
	github.com/slack-go/slack v0.12.5
	github.com/valyala/fasthttp v1.54.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/essentialkaos/depsy v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)
