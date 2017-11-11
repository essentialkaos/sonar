package svg

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

func GetPoint(color string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="12" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 12 16"><defs><path id="a" fill="` + color + `" d="M4.6 9.6Q4 10.2 4 11q0 .85.6 1.4.6.6 1.4.6.85 0 1.4-.6.6-.55.6-1.4 0-.8-.6-1.4Q6.85 9 6 9q-.8 0-1.4.6z"/></defs><use xlink:href="#a"/></svg>`
}

func GetAirplane() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 16 16"><defs><path id="a" fill="#60ADD8" d="M12.95 4.95q.25-.25.4-.6.3-.4.45-.85.15-.45.25-.85.05-.35-.15-.55-.15-.15-.5-.1-.45.05-.9.2-.4.2-.85.45-.3.2-.55.45L9.45 4.75 3.2 3.8H3l-.15.15L2.8 4l-.15.15v.2q0 .1.05.15.05.15.15.15L7.1 7.1 4.5 9.7l-1.9-.3q-.15-.05-.2 0-.15.05-.2.1l-.05.05q-.05.05-.1.2-.05.05 0 .2.05.05.05.15l.15.15L4.6 11.4l1.25 2.35q.05.15.1.2.1 0 .15.05.15.05.2 0 .15-.05.2-.1l.05-.05q.05-.05.1-.2.05-.05 0-.2l-.3-1.9 2.6-2.6 2.45 4.25q0 .1.15.15.05.05.15.05h.2l.15-.15.05-.05.15-.15v-.2L11.3 6.6l1.65-1.65z"/></defs><use xlink:href="#a"/></svg>`
}
