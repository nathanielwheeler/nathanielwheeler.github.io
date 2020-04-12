package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	res.Header().Set("Content-Type", "text/html")
	switch path {
	case "/":
		fmt.Fprint(res, "<h1>Nathaniel Wheeler</h1>")
	case "/contact":
		fmt.Fprint(res, "To get in touch, please send an email "+
		" to <a href=\"mailto:contact@nathanielwheeler.com\">"+
		"contact@nathanielwheeler.com</a>.")
	default:
		fmt.Fprint(res, "<h1>Page Not Found</h1>"+
		"<p>Please email me if you keep being sent to an invalid page.</p>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
