package marketsurge

import "net/http"

// Session carries caller-supplied browser cookies.
type Session struct {
	Cookies []*http.Cookie
}

// NewSession builds a Session from caller-supplied cookies.
func NewSession(cookies []*http.Cookie) Session {
	return Session{Cookies: cloneCookies(cookies)}
}

func cloneCookies(cookies []*http.Cookie) []*http.Cookie {
	if len(cookies) == 0 {
		return make([]*http.Cookie, 0)
	}

	cloned := make([]*http.Cookie, len(cookies))
	for i, cookie := range cookies {
		if cookie == nil {
			continue
		}

		copyCookie := *cookie // #nosec G124 - cloning caller-supplied cookies without changing security attributes.
		if len(cookie.Unparsed) > 0 {
			copyCookie.Unparsed = append([]string(nil), cookie.Unparsed...)
		}
		cloned[i] = &copyCookie
	}

	return cloned
}
