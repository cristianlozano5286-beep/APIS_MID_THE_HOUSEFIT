package main

import (
	"fmt"
	"log"
	"os"

	"api_mid_the_housefit/services"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Uso: go run ./cmd/gmail-test <correo-destino>")
	}

	servicio, err := services.NuevoServicioGmail()
	if err != nil {
		log.Fatal("No se pudo abrir el servicio OAuth de Gmail: ", err)
	}

	_, err = servicio.EnviarCorreo(
		os.Args[1],
		"Prueba de correo - The House Fit",
		"<p>Este es un correo de prueba de The House Fit. La integración con Gmail funciona correctamente.</p>",
		nil,
	)
	if err != nil {
		log.Fatal("No se pudo enviar el correo de prueba: ", err)
	}

	fmt.Println("Correo de prueba enviado correctamente a", os.Args[1])
}
