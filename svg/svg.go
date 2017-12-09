package svg

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// GetPoint return point svg with given color
func GetPoint(color string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none"><defs><path id="a" fill="` + color + `" d="M11 11q0-1.25-.85-2.15Q9.25 8 8 8q-1.2 0-2.1.85Q5 9.75 5 11t.9 2.1q.9.9 2.1.9 1.25 0 2.15-.9.85-.85.85-2.1z"/></defs><use xlink:href="#a"/></svg>`
}

// GetCircle return circle svg
func GetCircle() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#BBB" d="M2.9 7.9Q2 8.65 2 10q0 1.3.9 2.15.95.85 2.1.85 1.25 0 2.1-.85.9-.9.9-2.15t-.9-2.1Q6.3 7 5 7q-1.2 0-2.1.9m.6 2.1q0-.55.4-1.05.65-.45 1.1-.45.55 0 1.1.45.4.5.4 1.05 0 .6-.4 1.1-.5.4-1.1.4-.55 0-1.1-.4-.4-.5-.4-1.1z" id="a"/></defs><use xlink:href="#a"/></svg>`
}

// GetAirplane return airplane icon svg
func GetAirplane() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#589ABF" d="M14.35 2.5q-.4-.4-1.1-.5-.65-.05-1.4.3-.75.35-1.35.95l-1.7 1.7L4 3.55h-.5q-.2.1-.4.3l-1 1q-.3.3-.3.8.1.4.45.75L5.3 8.45l-1 1-1.15-.25q-.25-.05-.45.05-.3 0-.5.2l-.9.9q-.3.3-.3.8.05.45.5.7l2.1 1.4 1.4 2.1q.25.45.7.5.5 0 .8-.3l.9-.9q.2-.2.25-.45.05-.25 0-.5l-.25-1.15 1-1 2.05 3.05q.35.35.75.45.5 0 .8-.3l1-1q.2-.2.3-.4v-.5l-1.4-4.8 1.7-1.7q.6-.6.95-1.35.35-.75.3-1.4-.1-.7-.5-1.1m-2.9 1.7q.55-.55 1.15-.75.6-.3.85-.05t-.05.9q-.15.6-.7 1.15l-2.3 2.3L11.95 13l-.55.55L8.65 9.5 6 12.15l.35 1.65-.5.5-1.3-2-2-1.3.5-.5 1.7.4L7.4 8.25l-4.1-2.8.55-.55 5.3 1.6 2.3-2.3z" id="a"/></defs><use xlink:href="#a"/></svg>`
}

// GetAirplane return phone icon svg
func GetPhone() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#C67455" d="M11.45 2.05q-.45-.15-.85.05-.4.25-.6.65l-1.2 2.8q-.15.35-.05.75.15.45.45.7l1.1.9q-1.2 2.2-3.4 3.4L6 10.2q-.25-.3-.7-.45-.4-.1-.75.05L1.75 11q-.4.2-.65.6-.2.4-.05.85l.6 2.6q.1.4.4.7.35.25.8.25 3.3 0 6.15-1.65 2.7-1.6 4.35-4.35Q15 7.15 15 3.85q0-.45-.25-.8-.3-.3-.7-.4l-2.6-.6M11.2 3.3l2.5.6q0 2.9-1.5 5.45-1.45 2.4-3.85 3.85-2.55 1.5-5.45 1.5l-.6-2.5L5 11.05l1.5 1.85q2.05-.9 3.25-2.1Q11 9.55 11.9 7.5L10.05 6l1.15-2.7z" id="a"/></defs><use xlink:href="#a"/></svg>`
}
