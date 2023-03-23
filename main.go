package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/crontab"
	"github.com/joho/godotenv"
)

type Comprobante struct {
	IdComprobante     string `json:"idComprobante"`
	Fecha             string `json:"fecha"`
	Hora              string `json:"hora"`
	Nombre            string `json:"nombre"`
	CodigoComprobante string `json:"codigoComprobante"`
	Serie             string `json:"serie"`
	Numeracion        string `json:"numeracion"`
	Documento         string `json:"documento"`
	CodigoDocumento   string `json:"codigoDocumento"`
	NumeroDocumento   string `json:"numeroDocumento"`
	Informacion       string `json:"informacion"`
	Estado            string `json:"estado"`
	Tipo              string `json:"tipo"`
	Simbolo           string `json:"simbolo"`
	Total             string `json:"total"`
	Xmlsunat          string `json:"xmlsunat"`
	Xmldescripcion    string `json:"xmldescripcion"`
}

type Result struct {
	State       bool   `json:"state"`
	Accept      bool   `json:"accept"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Error struct {
	Message string `json:"message"`
}

type MessageWhatApp struct {
	Comprobante string
	Estado      string
	Descripcion string
}

func getCpeAllDocuments() {
    var url_server string = os.Getenv("URL_SERVER");

	var comprobantes []Comprobante
	var msgWhatApp []MessageWhatApp

	res, err := http.Get(url_server+"/app/controller/VentaController.php?type=listComprobantesExternal")
	if err != nil {
		sendWhatsAppMessageFail("No se puede establecer una conexión ya que el equipo de destino denegó expresamente dicha conexión.")
		// sendWhatsAppMessageFail(strings.TrimSpace(string(err.Error())))
	} else {
		defer res.Body.Close()

		if res.StatusCode == 200 {
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal([]byte(body), &comprobantes)

			if len(comprobantes) > 0 {
				msgWhatApp = sendAllDocuments(comprobantes)
				var cadena string
				for _, obj := range msgWhatApp {
					cadena += fmt.Sprintf("   ---> %s %s %s", obj.Comprobante, obj.Estado, obj.Descripcion)
				}
				sendWhatsAppMessage(cadena)
			}
		}
	}

}

func sendAllDocuments(comprobantes []Comprobante) []MessageWhatApp {

	var msgWA MessageWhatApp
	var arrayMsgWA []MessageWhatApp

	for _, comp := range comprobantes {

		if comp.Tipo == "v" {
			arrayMsgWA = append(arrayMsgWA, validateBolateFactura(comp, msgWA))
		} else if comp.Tipo == "gui" {
			arrayMsgWA = append(arrayMsgWA, validateGuia(comp, msgWA))
		} else if comp.Tipo == "nc" {
			arrayMsgWA = append(arrayMsgWA, validateNotaCredito(comp, msgWA))
		}
	}

	return arrayMsgWA
}

func validateBolateFactura(comp Comprobante, msgWA MessageWhatApp) MessageWhatApp {

	if comp.Xmlsunat == "0" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante ya fue aceptado"

	} else if comp.Xmlsunat == "1033" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante aceptado ya no se puede declarar mas de 2 veces"

	} else if len(comp.Serie) != 4 {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "La serie no cumple con el número de caracteres establecido"

	} else if comp.CodigoComprobante == "03" { // CODIGO DE BOLETA = 03

		if !strings.HasPrefix(strings.ToUpper(comp.Serie), "B") {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "La serie no cumple con el formato de facturación"

		} else if len(comp.NumeroDocumento) != 8 {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El número de documento no cumple con el número de caracteres establecido"

		} else if comp.CodigoDocumento == "1" || comp.CodigoDocumento == "0" {
			// ENVIO DE BOLETA
			msgWA = automaticDeliveryVouchers(comp, msgWA)

		} else {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El codigo de documento es incorrecto"

		}

	} else if comp.CodigoComprobante == "01" { // Codigo de Factura = 01

		if !strings.HasPrefix(strings.ToUpper(comp.Serie), "F") {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "La serie no cumple con el formato de facturación"

		} else if len(comp.NumeroDocumento) != 11 {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El número de documento no cumple con el número de caracteres establecido"

		} else if comp.CodigoDocumento == "6" {
			// ENVIO DE FACTURA
			msgWA = automaticDeliveryVouchers(comp, msgWA)

		} else {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El codigo de documento es incorrecto"
		}

	} else {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "El codigo de comprobante es incorrecto"

	}

	return msgWA
}

func validateGuia(comp Comprobante, msgWA MessageWhatApp) MessageWhatApp {

	if comp.Xmlsunat == "0" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante ya fue enviado"

	} else if comp.Xmlsunat == "1033" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante aceptado ya no se puede declarar mas de 2 veces"

	} else if len(comp.Serie) != 4 {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "La serie no cumple con el número de caracteres establecido"

	} else if comp.CodigoComprobante == "09" {

		if !strings.HasPrefix(strings.ToUpper(comp.Serie), "T") {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "La serie no cumple con el formato de facturación"

		} else if len(comp.NumeroDocumento) != 11 && len(comp.NumeroDocumento) != 8 {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El número de documento no cumple con el número de caracteres establecido"

		} else if comp.CodigoDocumento == "6" || comp.CodigoDocumento == "1" || comp.CodigoDocumento == "0" {
			msgWA = automaticDeliveryVouchers(comp, msgWA)

		} else {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El codigo de documento es incorrecto"

		}

	} else {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "El codigo de comprobante es incorrecto"

	}

	return msgWA
}

func validateNotaCredito(comp Comprobante, msgWA MessageWhatApp) MessageWhatApp {

	if comp.Xmlsunat == "0" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante ya fue enviado"

	} else if comp.Xmlsunat == "1033" {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "COMPROBANTE ACEPTADO"
		msgWA.Descripcion = "El comprobante aceptado ya no se puede declarar mas de 2 veces"

	} else if len(comp.Serie) != 4 {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "La serie no cumple con el número de caracteres establecido"

	} else if comp.CodigoComprobante == "07" {

		if !strings.HasPrefix(strings.ToUpper(comp.Serie), "BN") && !strings.HasPrefix(strings.ToUpper(comp.Serie), "FN") {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "La serie no cumple con el formato de facturación"

		} else if len(comp.NumeroDocumento) != 11 && len(comp.NumeroDocumento) != 8 {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El número de documento no cumple con el número de caracteres establecido"

		} else if comp.CodigoDocumento == "6" || comp.CodigoDocumento == "1" || comp.CodigoDocumento == "0" {
			msgWA = automaticDeliveryVouchers(comp, msgWA)

		} else {
			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "CORREGIR"
			msgWA.Descripcion = "El codigo de documento es incorrecto"

		}

	} else {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "CORREGIR"
		msgWA.Descripcion = "El codigo de comprobante es incorrecto"

	}

	return msgWA
}

func automaticDeliveryVouchers(comp Comprobante, msgWA MessageWhatApp) MessageWhatApp {

	var result Result
	var error Error
	var url string
	var url_server string = os.Getenv("URL_SERVER");

	if comp.Tipo == "v" {
		url = url_server+"/app/examples/boleta.php?idventa="
	} else if comp.Tipo == "gui" {
		url = url_server+"/app/examples/guiaremision.php?idGuiaRemision="
	} else if comp.Tipo == "nc" {
		url = url_server+"/app/examples/notacredito.php?idNotaCredito="
	}

	res, err := http.Get(url + comp.IdComprobante)
	if err != nil {
		msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
		msgWA.Estado = "ERROR"
		msgWA.Descripcion = "No se pudo conectar al servidor para el envio de comprobantes"
	} else {
		defer res.Body.Close()

		if res.StatusCode == 200 {
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal([]byte(body), &result)

			if result.State && result.Accept {
				msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
				msgWA.Estado = "ENVIADO"
				msgWA.Descripcion = result.Description
			} else {
				msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
				msgWA.Estado = "CORREGIR"
				msgWA.Descripcion = result.Description
			}
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal([]byte(body), &error)

			msgWA.Comprobante = comp.Serie + "-" + comp.Numeracion
			msgWA.Estado = "ADVERTENCIA"
			msgWA.Descripcion = error.Message
		}
	}

	return msgWA
}

func sendWhatsAppMessage(msgWA string) {

	url := "https://graph.facebook.com/v16.0/114181978288370/messages"

	message := `{
		"messaging_product": "whatsapp",
		"preview_url": false,
		"recipient_type": "individual",
		"to": "51931341082",
		"type": "text",
		"text": {
			"body": "%s"
		}
	}`

	payload := strings.NewReader(fmt.Sprintf(message, msgWA))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer EAAK4CE2kFXABABuxeQeEjEJFWe08nHlw0RDIEyIX2SyJjvTW0NPiQAdCaKxOkKjZCOQ0ALI6iDrTZBa8WXoT5kA1WjZBDTa6lkHZCqcHjOgghxQ8RDy9jWUm6thtLSYauG2nxS16HOgGOSVDpxSwxjxxPo9ujVZC2uaDP9P2wkHzvCSHw3Am6go3aoesQ6P8wjbQjehmjiQZDZD")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			fmt.Println(res.Status)
		} else {
			// fmt.Printf("%s", msgWA)
			body, _ := ioutil.ReadAll(res.Body)
			message := strings.TrimSpace(string(body))
			fmt.Println(message)
		}
	}

}

func sendWhatsAppMessageFail(msg string) {
	url := "https://graph.facebook.com/v16.0/114181978288370/messages"

	message := `{
		"messaging_product": "whatsapp",
		"preview_url": false,
		"recipient_type": "individual",
		"to": "51931341082",
		"type": "text",
		"text": {
			"body": "%s"
		}
	}`

	payload := strings.NewReader(fmt.Sprintf(message, msg))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer EAAK4CE2kFXABABuxeQeEjEJFWe08nHlw0RDIEyIX2SyJjvTW0NPiQAdCaKxOkKjZCOQ0ALI6iDrTZBa8WXoT5kA1WjZBDTa6lkHZCqcHjOgghxQ8RDy9jWUm6thtLSYauG2nxS16HOgGOSVDpxSwxjxxPo9ujVZC2uaDP9P2wkHzvCSHw3Am6go3aoesQ6P8wjbQjehmjiQZDZD")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error de envio de mensaje por whatsapp")
	} else {
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			fmt.Println(res.Status)
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			message := strings.TrimSpace(string(body))
			fmt.Println(message)
		}
	}
}

func main() {

	time.LoadLocation("America/Lima")

	godotenv.Load();

	var go_port string = os.Getenv("GO_PORT");

	ctab := crontab.New()

	err := ctab.AddJob("53 08 * * *", func() {		
		getCpeAllDocuments()
	})

	if err != nil {
		log.Fatal(err)
	}

	// ctab.MustAddJob("* * * * *", myFunc)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// getCpeAllDocuments()
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(go_port)
}

// func task() {
// 	b := []byte("Hola mundo!\n")
// 	dt := time.Now()
// 	// fob := dt.Format("15:04:05")
// 	formattime := strconv.Itoa(dt.Hour()) + "-" + strconv.Itoa(dt.Minute()) + "-" + strconv.Itoa(dt.Second())
// 	err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

