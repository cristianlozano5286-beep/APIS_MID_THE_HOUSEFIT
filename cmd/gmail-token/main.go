package main

import (
	"fmt"
	"log"

	"api_mid_the_housefit/services"
)

func main() {
	if _, err := services.NuevoServicioGmail(); err != nil {
		log.Fatal("No se pudo completar la autorización OAuth de Gmail: ", err)
	}

	fmt.Println("Autorización completada. El token está guardado en token/token.json.")
}
