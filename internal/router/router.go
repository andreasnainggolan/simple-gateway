package router

import "strings"

/*
ProtectConfig menyimpan informasi proteksi per-route.
V1: hanya API key.
Nanti bisa berkembang (rate limit, basic auth, dll).
*/
type ProtectConfig struct {
	APIKey bool
}

/*
Route merepresentasikan satu entry dari gateway.yaml
yang sudah dikompilasi.
*/
type Route struct {
	Host      string
	PathRaw   string
	Path      PathPattern
	ForwardTo string
	Protect   *ProtectConfig
}

/*
MatchResult adalah hasil pencocokan router.
*/
type MatchResult struct {
	Route  Route
	Params map[string]string
}

/*
Router menyimpan seluruh route dan menyediakan fungsi Match().
*/
type Router struct {
	routes []Route
}

/*
New membuat instance Router.
*/
func New(routes []Route) *Router {
	return &Router{routes: routes}
}

/*
Match mencari route pertama yang cocok dengan requestHost dan requestPath.

Aturan host:
- Jika route.Host kosong → match semua host
- Jika route.Host diisi → harus sama persis
*/
func (r *Router) Match(requestHost, requestPath string) (MatchResult, bool) {
	requestHost = strings.TrimSpace(strings.ToLower(requestHost))

	for _, rt := range r.routes {
		if rt.Host != "" && strings.ToLower(rt.Host) != requestHost {
			continue
		}

		ok, params := MatchPath(rt.Path, requestPath)
		if ok {
			return MatchResult{
				Route:  rt,
				Params: params,
			}, true
		}
	}

	return MatchResult{}, false
}
