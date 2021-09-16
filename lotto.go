package main

import (
	"lotto/lottoapi"
	"lotto/lottohttp"
	"net/http"
)

func main() {

	lottoapi.InitNutzer()

	http.HandleFunc("/", lottohttp.HttpResponder)
	http.ListenAndServe(":8080", nil)

}
