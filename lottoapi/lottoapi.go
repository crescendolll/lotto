package lottoapi

import (
	"fmt"
	"lotto/database"
	"lotto/lottojson"
	"lotto/lottologic"
)

var aktiveNutzer map[string]lottologic.Nutzer

// Ausgabe ist interface{} - beliebig, weil jede Response erzeugt werden kann
func HandleRequest(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	switch apirequest.Methode {
	case "login":
		response = CreateLoginResponse(apirequest)
	}

	return response

}

func CreateLoginResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	if apirequest.Param["name"] == "" || apirequest.Param["pwhash"] == "" {
		response = lottojson.Errorresponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		databasehandle := database.OpenLottoConnection()
		nutzer, err := database.SelectFromSpielerByName(databasehandle, apirequest.Param["name"])
		database.CloseLottoConnection(databasehandle)

		if err != nil {
			response = lottojson.Errorresponse{
				Errormessage: err.Error(),
			}
		} else {
			if nutzer.Pw_hash != apirequest.Param["pwhash"] {
				response = lottojson.Errorresponse{
					Errormessage: "Passwort oder Benutzername sind nicht korrekt",
				}
			} else {
				auth := LoginNutzer(nutzer)
				response = lottojson.LoginResponse{
					Errormessage: "",
					IstSpieler:   nutzer.Ist_spieler,
					Auth:         auth,
				}
			}
		}
	}

	for key, value := range aktiveNutzer {
		fmt.Printf("Nutzer %s mit PW %s unter Token %s aktiv\n", value.Benutzername, value.Pw_hash, key)
	}

	return response

}

func InitNutzer() {

	aktiveNutzer = make(map[string]lottologic.Nutzer)

}

func LoginNutzer(nutzer lottologic.Nutzer) string {

	auth := nutzer.Benutzername + ":" + nutzer.Pw_hash

	aktiveNutzer[auth] = nutzer

	return auth
}
