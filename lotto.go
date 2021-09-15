package main

import (
	"lotto/database"
	"lotto/lottohttp"
	"net/http"
)

func main() {

	databasehandle := database.OpenLottoConnection()

	database.CloseLottoConnection(databasehandle)

	http.HandleFunc("/", lottohttp.HttpResponder)
	http.ListenAndServe(":8080", nil)

}
