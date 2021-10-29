package lottojson

import "lotto/lottologic"

type LottoRequest struct {
	Methode string            `json:"methode"`
	Param   map[string]string `json:"param"`
	Auth    string            `json:"auth"`
}

type ErrorResponse struct {
	Errormessage string `json:"errormessage"`
}

type LoginResponse struct {
	Errormessage string `json:"errormessage"`
	IstSpieler   bool   `json:"istspieler"`
	Auth         string `json:"auth"`
}

type RegistrationResponse struct {
	Errormessage string `json:"errormessage"`
	Auth         string `json:"auth"`
}

type GetZiehungenResponse struct {
	Errormessage string   `json:"errormessage"`
	Ziehungen    []string `json:"ziehungstage"`
}

type GetTippsResponse struct {
	Errormessage     string                      `json:"errormessage"`
	Tippauszahlungen []lottologic.Tippauszahlung `json:"statistik"`
}

type GetZiehungenMitAuszahlungenResponse struct {
	Errormessage         string                         `json:"errormessage"`
	Ziehungsauszahlungen []lottologic.Ziehungsstatistik `json:"statistik"`
}
