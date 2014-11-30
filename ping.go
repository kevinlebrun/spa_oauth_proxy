package main

import "net/http"

type PingHandler struct {
	Response string
}

func (p PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(p.Response))
}
