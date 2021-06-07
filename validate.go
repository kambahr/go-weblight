package main

import (
	"net/http"
	"strings"
)

func (ws *website) validateRequest(w http.ResponseWriter, r *http.Request) bool {

	// For performance reasons, keep all logic here (no other functions),
	// even for short lookups that may repeat a couple of times!

	// Note: blocked ips have already been dropped.

	rPath := strings.ToLower(r.URL.Path)

	// Method allowed
	failed := true
	s := r.Method
	for i := 0; i < len(ws.Config.HTTP.AllowedMethods); i++ {
		if s == ws.Config.HTTP.AllowedMethods[i] {
			failed = false
			break
		}
	}
	if failed {
		ws.displayHTTPErrMsg(w, 405)
		return false
	}

	// This is the order: restrict-paths, exclude-path, forward-paths

	// restrict-paths
	for i := 0; i < len(ws.Config.URLPaths.Restrict); i++ {
		sl := strings.ToLower(ws.Config.URLPaths.Restrict[i])
		// On-error  the err-text is placed instead of the
		// expected value starting with: ~@error
		if strings.HasPrefix(sl, "~@error") {
			continue
		}
		if sl == rPath {
			ws.displayHTTPErrMsg(w, http.StatusUnauthorized)
			return false
		}
	}

	// exclude-path
	for i := 0; i < len(ws.Config.URLPaths.Exclude); i++ {
		sl := strings.ToLower(ws.Config.URLPaths.Exclude[i])
		if strings.HasPrefix(sl, "~@error") {
			continue
		}
		if sl == rPath {
			ws.displayHTTPErrMsg(w, http.StatusNotFound)
			return false
		}
	}

	// forward-paths
	for i := 0; i < len(ws.Config.URLPaths.Forward); i++ {
		sl := strings.ToLower(ws.Config.URLPaths.Forward[i])
		v := strings.Split(sl, "|")
		left := ""
		right := ""
		if len(v) > 1 {
			left = v[0]
			right = v[1]
		}
		if right == "" || strings.HasPrefix(right, "~@error") {
			continue
		}
		if left == rPath {
			http.Redirect(w, r, right, http.StatusTemporaryRedirect)
			return false
		}
	}

	return true
}
