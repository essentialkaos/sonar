package svg

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// GetPoint return point svg with given color
func GetPoint(color string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="12" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 16 16"><defs><path id="a" fill="` + color + `" d="M11 11q0-1.25-.85-2.15Q9.25 8 8 8q-1.2 0-2.1.85Q5 9.75 5 11t.9 2.1q.9.9 2.1.9 1.25 0 2.15-.9.85-.85.85-2.1z"/></defs><use xlink:href="#a"/></svg>`
}

// GetCircle return circle svg
func GetCircle() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 16 16"><defs><path id="a" fill="#BBB" d="M11 11q0-1.25-.85-2.15Q9.25 8 8 8q-1.2 0-2.1.85Q5 9.75 5 11t.9 2.1q.9.9 2.1.9 1.25 0 2.15-.9.85-.85.85-2.1M6.95 9.95Q7.4 9.5 8 9.5q.65 0 1.1.45.4.45.4 1.05 0 .65-.4 1.05-.45.45-1.1.45-.6 0-1.05-.45-.45-.4-.45-1.05 0-.6.45-1.05z"/></defs><use xlink:href="#a"/></svg>`
}

// GetAirplane return airplane icon svg
func GetAirplane() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 16 16"><defs><path id="a" fill="#60ADD8" d="M12.5 4.15q-.4.2-.85.45-.3.2-.55.45L9.45 6.7 3.2 5.75H3l-.15.15-.05.05-.15.15v.2q0 .1.05.15.05.15.15.15L7.1 9.05l-2.6 2.6-1.9-.3q-.15-.05-.2 0-.15.05-.2.1l-.05.05q-.05.05-.1.2-.05.05 0 .2.05.05.05.15l.15.15 2.35 1.15 1.25 2.35q.05.15.1.2.1 0 .15.05.15.05.2 0 .15-.05.2-.1l.05-.05q.05-.05.1-.2.05-.05 0-.2l-.3-1.9 2.6-2.6 2.45 4.25q0 .1.15.15.05.05.15.05h.2l.15-.15.05-.05.15-.15v-.2l-.95-6.25 1.65-1.65q.25-.25.4-.6.3-.4.45-.85.15-.45.25-.85.05-.35-.15-.55-.15-.15-.5-.1-.45.05-.9.2z"/></defs><use xlink:href="#a"/></svg>`
}
