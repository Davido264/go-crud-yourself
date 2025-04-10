package main

import "net/http"

func hasJsonBody(req *http.Request) bool {
	return req.Header.Get("Content-Type") == "application/json"
}
