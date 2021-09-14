package main

import "lotto/database"

func main() {

	databasehandle := database.OpenLottoConnection()

	database.CloseLottoConnection(databasehandle)

}
