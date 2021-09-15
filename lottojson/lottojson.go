package lottojson

type LottoRequest struct {
	Methode string            `json:"methode"`
	Param   map[string]string `json:"param"`
}

type LoginResponse struct {
	Errormessage string `json:"errormessage"`
	IstSpieler   bool   `json:"istspieler"`
	Auth         string `json:"auth"`
}
