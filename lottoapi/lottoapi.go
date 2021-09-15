package lottoapi

import "lotto/lottojson"

func HandleRequest(apirequest lottojson.LottoRequest) lottojson.LoginResponse {

	var response lottojson.LoginResponse

	response = lottojson.LoginResponse{
		Errormessage: "",
		IstSpieler:   true,
		Auth:         apirequest.Param["name"] + ":" + apirequest.Param["pwhash"],
	}

	return response

}
