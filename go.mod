module github.com/essentialkaos/sonar

go 1.19

replace github.com/slack-go/slack v0.12.2 => github.com/essentialkaos/slack v0.0.0-20230418134115-ed5fe51c0353

require (
	github.com/essentialkaos/depsy v1.1.0
	github.com/essentialkaos/ek/v12 v12.91.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/slack-go/slack v0.12.2
	github.com/valyala/fasthttp v1.51.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
