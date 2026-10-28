package mw

import "net/http"

var publicRoutes = map[string]struct{}{
	http.MethodGet + " /api/auth/getCaptcha":           {},
	http.MethodPost + " /api/auth/login":               {},
	http.MethodPost + " /api/auth/refresh":             {},
	http.MethodPost + " /api/auth/verification/create": {},
	http.MethodPost + " /api/auth/verification/verify": {},
	http.MethodGet + " /api/checkin/template":          {},
	http.MethodGet + " /api/checkin/state":             {},
	http.MethodGet + " /api/checkin/qrcode":            {},
	http.MethodGet + " /api/ws":                        {},
}

func isPublicRoute(method, path string) bool {
	_, exists := publicRoutes[method+" "+path]
	return exists
}
