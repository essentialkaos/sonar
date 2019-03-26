package svg

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// GetBullet return bullet svg with given color
func GetBullet(color string) string {
	switch color {
	case "":
		return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#fff" fill-opacity="0" d="M16 0H0v16h16V0z" id="a"/></defs><use xlink:href="#a"/></svg>`
	default:
		return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="` + color + `" d="M5.85 13.1q.9.9 2.15.9 1.3 0 2.1-.9.9-.85.9-2.1t-.9-2.15Q9.3 8 8 8q-1.25 0-2.15.85Q5 9.75 5 11t.85 2.1z" id="a"/></defs><use xlink:href="#a"/></svg>`
	}
}

// GetDND return circle with Z
func GetDND() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#BBB" d="M12.45 7.55V6.6H8.8v1.15h1.65l-1.8 2.1v.75h3.85V9.45h-1.65l1.6-1.9M9 8.15Q8.56074219 8 8 8q-1.25 0-2.15.85Q5 9.75 5 11t.85 2.1q.9.9 2.15.9 1.3 0 2.1-.9.77109375-.728125.85-1.75h-1.5q-.07832031.4140625-.4.7-.4.45-1.05.45-.6 0-1.05-.45-.45-.4-.45-1.05 0-.6.45-1.05.39277344-.39277344.9-.45L9 8.15z" id="a"/></defs><use xlink:href="#a"/></svg>`
}

// GetCircle return circle svg
func GetCircle() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="16" height="16"><defs><path fill="#BBB" d="M8 8q-1.25 0-2.15.85Q5 9.75 5 11t.85 2.1q.9.9 2.15.9 1.3 0 2.1-.9.9-.85.9-2.1t-.9-2.15Q9.3 8 8 8m-1.05 4.05q-.45-.4-.45-1.05 0-.6.45-1.05Q7.4 9.5 8 9.5q.65 0 1.05.45.45.45.45 1.05 0 .65-.45 1.05-.4.45-1.05.45-.6 0-1.05-.45z" id="a"/></defs><use xlink:href="#a"/></svg>`
}

// GetAirplane return airplane icon svg
func GetAirplane() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="20" height="16"><defs><path fill="#589ABF" d="M14.35 2.5q-.4-.4-1.1-.5-.65-.05-1.4.3-.75.35-1.35.95l-1.7 1.7L4 3.55h-.5q-.2.1-.4.3l-1 1q-.3.3-.3.8.1.4.45.75L5.3 8.45l-1 1-1.15-.25q-.25-.05-.45.05-.3 0-.5.2l-.9.9q-.3.3-.3.8.05.45.5.7l2.1 1.4 1.4 2.1q.25.45.7.5.5 0 .8-.3l.9-.9q.2-.2.25-.45.05-.25 0-.5l-.25-1.15 1-1 2.05 3.05q.35.35.75.45.5 0 .8-.3l1-1q.2-.2.3-.4v-.5l-1.4-4.8 1.7-1.7q.6-.6.95-1.35.35-.75.3-1.4-.1-.7-.5-1.1m-2.9 1.7q.55-.55 1.15-.75.6-.3.85-.05t-.05.9q-.15.6-.7 1.15l-2.3 2.3L11.95 13l-.55.55L8.65 9.5 6 12.15l.35 1.65-.5.5-1.3-2-2-1.3.5-.5 1.7.4L7.4 8.25l-4.1-2.8.55-.55 5.3 1.6 2.3-2.3z" id="a"/></defs><use xlink:href="#a" transform="translate(4)"/></svg>`
}

// GetAirplane return phone icon svg
func GetPhone() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="none" width="20" height="16"><defs><path fill="#C67455" d="M14.75 3.05q-.3-.3-.7-.4l-2.6-.6q-.45-.15-.85.05-.4.25-.6.65l-1.2 2.8q-.15.35-.05.75.15.45.45.7l1.1.9q-1.2 2.2-3.4 3.4L6 10.2q-.25-.3-.7-.45-.4-.1-.75.05L1.75 11q-.4.2-.65.6-.2.4-.05.85l.6 2.6q.1.4.4.7.35.25.8.25 3.3 0 6.15-1.65 2.7-1.6 4.35-4.35Q15 7.15 15 3.85q0-.45-.25-.8m-1.05.85q0 2.9-1.5 5.45-1.45 2.4-3.85 3.85-2.55 1.5-5.45 1.5l-.6-2.5L5 11.05l1.5 1.85q2.05-.9 3.25-2.1Q11 9.55 11.9 7.5L10.05 6l1.15-2.7 2.5.6z" id="a"/></defs><use xlink:href="#a" transform="translate(4)"/></svg>`
}
