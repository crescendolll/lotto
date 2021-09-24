package lottohttp

import (
	"encoding/json"
	"fmt"
	"lotto/lottoapi"
	"lotto/lottojson"
	"net/http"
)

func OpenLottoServer() {

	fmt.Println("Lottoserver läuft")

	lottoapi.InitNutzer()

	http.HandleFunc("/", HttpResponder)
	http.ListenAndServe(":8080", nil)

}

func HttpResponder(responsewriter http.ResponseWriter, request *http.Request) {

	responsewriter.Header().Set("Content-Type", "application/json")

	headerContentType := request.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		apiresponse := lottojson.ErrorResponse{
			Errormessage: "not a JSON",
		}
		jsonResp, _ := json.Marshal(apiresponse)
		responsewriter.Write(jsonResp)
		return
	}
	var apirequest lottojson.LottoRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&apirequest)
	if err != nil {
		apiresponse := lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
		jsonResp, _ := json.Marshal(apiresponse)
		responsewriter.Write(jsonResp)
		return
	}

	apiresponse := lottoapi.HandleRequest(apirequest)

	jsonResp, _ := json.Marshal(apiresponse)
	responsewriter.Write(jsonResp)

	return
}
