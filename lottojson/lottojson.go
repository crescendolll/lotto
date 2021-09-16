package lottojson

type LottoRequest struct {
	Methode string            `json:"methode"`
	Param   map[string]string `json:"param"`
}

type Errorresponse struct {
	Errormessage string `json:"errormessage"`
}

type LoginResponse struct {
	Errormessage string `json:"errormessage"`
	IstSpieler   bool   `json:"istspieler"`
	Auth         string `json:"auth"`
}

type LogoutResponse struct {
	Errormessage string `json:"errormessage"`
	Auth         string `json:"auth"`
}
