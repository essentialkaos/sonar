package svg

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// GetPoint return point svg with given color
func GetPoint(color string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="12" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 12 16"><defs><path id="a" fill="` + color + `" d="M6 9q-1.2 0-2 1-1 .8-1 2 0 1.3 1 2 .8 1 2 1 1.3 0 2-1 1-.8 1-2 0-1-1-2-.8-1-2-1z"/></defs><use xlink:href="#a"/></svg>`
}

// GetCircle return circle svg
func GetCircle() string {
	return `
388 bytes - 45.51% saving
<svg xmlns="http://www.w3.org/2000/svg" width="12" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 12 16"><defs><path id="a" fill="#B2B2B2" d="M4 10q-1 .8-1 2 0 1.3 1 2v.2q.8.8 2 .8 1.3 0 2-.8h.2Q9 13 9 12q0-1.2-.8-2Q7 9 6 9q-1.2 0-2 1m3 1q.3.5.3 1T7 13q-.5.3-1 .3T5 13q-.3-.5-.3-1t.4-1q.6-.3 1-.3t1 .4z"/></defs><use xlink:href="#a"/></svg>`
}

// GetAirplane return airplane icon svg
func GetAirplane() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" viewBox="0 0 16 16"><defs><path id="a" fill="#60ADD8" d="M12.5 4.15q-.4.2-.85.45-.3.2-.55.45L9.45 6.7 3.2 5.75H3l-.15.15-.05.05-.15.15v.2q0 .1.05.15.05.15.15.15L7.1 9.05l-2.6 2.6-1.9-.3q-.15-.05-.2 0-.15.05-.2.1l-.05.05q-.05.05-.1.2-.05.05 0 .2.05.05.05.15l.15.15 2.35 1.15 1.25 2.35q.05.15.1.2.1 0 .15.05.15.05.2 0 .15-.05.2-.1l.05-.05q.05-.05.1-.2.05-.05 0-.2l-.3-1.9 2.6-2.6 2.45 4.25q0 .1.15.15.05.05.15.05h.2l.15-.15.05-.05.15-.15v-.2l-.95-6.25 1.65-1.65q.25-.25.4-.6.3-.4.45-.85.15-.45.25-.85.05-.35-.15-.55-.15-.15-.5-.1-.45.05-.9.2z"/></defs><use xlink:href="#a"/></svg>`
}
