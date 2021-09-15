package lottohttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"lotto/lottoapi"
	"lotto/lottojson"
	"net/http"
)

func HttpResponder(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		fmt.Println("not a JSON")
		return
	}

	var apirequest lottojson.LottoRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&apirequest)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			fmt.Println("Bad Request. Wrong Type provided for field " + unmarshalErr.Field)
		} else {
			fmt.Println("Bad Request " + err.Error())
		}
		return
	}

	apiresponse := lottoapi.HandleRequest(apirequest)

	jsonResp, _ := json.Marshal(apiresponse)
	w.Write(jsonResp)

	return
}
