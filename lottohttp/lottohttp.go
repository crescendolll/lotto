package lottohttp

import (
	"encoding/json"
	"lotto/lottoapi"
	"lotto/lottojson"
	"lotto/lottolog"
	"net/http"
)

func StarteLottoServer() {

	lottoapi.InitialisiereNutzerliste()

	http.HandleFunc("/", BeantworteHTTP)
	http.ListenAndServe(":8080", nil)

}

func BeantworteHTTP(responsewriter http.ResponseWriter, request *http.Request) {

	responsewriter.Header().Set("Content-Type", "application/json")

	headerContentType := request.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		apiresponse := lottojson.ErrorResponse{
			Errormessage: "Header Content Type nicht korrekt mit application/json gesetzt",
		}
		jsonResponse, jsonFehler := json.Marshal(apiresponse)
		if jsonFehler == nil {
			responsewriter.Write(jsonResponse)
		} else {
			lottolog.FehlerLogger.Println(jsonFehler.Error())
		}
		return
	}
	var apirequest lottojson.LottoRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	fehler := decoder.Decode(&apirequest)
	if fehler != nil {
		apiresponse := lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
		jsonResponse, jsonFehler := json.Marshal(apiresponse)
		if jsonFehler == nil {
			responsewriter.Write(jsonResponse)
		} else {
			lottolog.FehlerLogger.Println(jsonFehler.Error())
		}
		return
	}

	apiresponse := lottoapi.BearbeiteRequest(apirequest)

	jsonResponse, jsonFehler := json.Marshal(apiresponse)
	if jsonFehler == nil {
		responsewriter.Write(jsonResponse)
	} else {
		lottolog.FehlerLogger.Println(jsonFehler.Error())
	}

	return
}
