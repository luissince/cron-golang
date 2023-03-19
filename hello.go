package main

import (
	// "github.com/go-co-op/gocron"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/crontab"
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

type MessageWhatApp struct {
	Comprobante string
	Estado      string
	Descripcion string
}

func getCpeBoletaFactura() {

	var comprobantes []Comprobante
	// var msgWhatApp []MessageWhatApp

	res, err := http.Get("http://localhost:9000/app/controller/VentaController.php?type=listCpeBoletaFactura")
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal([]byte(body), &comprobantes)

		// msgWhatApp = automaticDeliveryVouchers(&comprobantes)
		// SendWhatsAppMessage(msgWhatApp)
		fmt.Printf("%+v", comprobantes)

	} else {
		// resulMessage = "Estado de peticóon diferente de 200"
		return
	}

	/*
		F001-1
		ENVIADO

		F002
		POR REVISAR
		DETALASDASDASD
	*/

}

func automaticDeliveryVouchers(comprobantes *[]Comprobante) []MessageWhatApp {

	var msgWA MessageWhatApp
	var arrayMsgWA []MessageWhatApp

	for _, comp := range *comprobantes {

		if comp.Tipo == "v" {
			if comp.Xmlsunat == "0" {
				return nil
			}

			res, err := http.Get("http://localhost:9000/app/examples/boleta.php?idventa=" + comp.IdComprobante)
			if err != nil {
				return nil
			}
			defer res.Body.Close()

			if res.StatusCode == 200 {

				body, _ := ioutil.ReadAll(res.Body)

				var result Result
				json.Unmarshal([]byte(body), &result)

				if result.State && result.Accept {

					msgWA.Comprobante = comp.Serie + " - " + comp.Numeracion
					msgWA.Estado = "ENVIADO"
					msgWA.Descripcion = result.Description

					// message = result.Description

					arrayMsgWA = append(arrayMsgWA, msgWA)

					// return
				}
			}

			// fmt.Println(res.Body)

			// if comp.Estado == "3" {
			// 	if comp.Xmlsunat == "1032" {

			// 		boletaFactura := strings.ToUpper(comp.Serie)
			// 		if strings.Contains(boletaFactura, "B") {
			// 			res, err := http.Get("http://localhost:9000/app/examples/resumen.php?idventa=" + comp.IdComprobante)
			// 			if err != nil {
			// 				fmt.Println("Ocurrio un error: ", err)
			// 			}
			// 			defer res.Body.Close()
			// 			fmt.Println(res.Body)
			// 		}

			// 		if strings.Contains(boletaFactura, "F") {
			// 			res, err := http.Get("http://localhost:9000/app/examples/comunicacionbaja.php?idventa=" + comp.IdComprobante)
			// 			if err != nil {
			// 				fmt.Println("Ocurrio un error: ", err)
			// 			}
			// 			defer res.Body.Close()
			// 			fmt.Println(res.Body)
			// 		}
			// 	}
			// }

			// return res.Status
		}

		// if comp.Tipo == "nc" {
		// 	if comp.Xmlsunat == "0" {
		// 		return
		// 	}

		// 	res, err := http.Get("http://localhost:9000/app/examples/notacredito.php?idNotaCredito=" + comp.IdComprobante)
		// 	if err != nil {
		// 		return fmt.Sprint("Ocurrio un error: ", err)
		// 	}
		// 	defer res.Body.Close()

		// 	if res.StatusCode == 200 {

		// 		body, _ := ioutil.ReadAll(res.Body)

		// 		var result Result
		// 		json.Unmarshal([]byte(body), &result)

		// 		if result.State && result.Accept {
		// 			message = result.Description
		// 			return message
		// 		}
		// 	}
		// }

		// if comp.Tipo == "gui" {
		// 	if comp.Xmlsunat == "0" {
		// 		return
		// 	}

		// 	res, err := http.Get("http://localhost:9000/app/examples/guiaremision.php?idGuiaRemision=" + comp.IdComprobante)
		// 	if err != nil {
		// 		return fmt.Sprint("Ocurrio un error: ", err)
		// 	}
		// 	defer res.Body.Close()

		// 	if res.StatusCode == 200 {

		// 		body, _ := ioutil.ReadAll(res.Body)

		// 		var result Result
		// 		json.Unmarshal([]byte(body), &result)

		// 		if result.State && result.Accept {
		// 			message = result.Description
		// 			return message
		// 		}
		// 	}
		// }

	}

	return arrayMsgWA
}

func SendWhatsAppMessage(msgWA []MessageWhatApp) {

	url := "https://graph.facebook.com/v15.0/114181978288370/messages"

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
	// req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")
	req.Header.Add("Authorization", "Bearer EAAK4CE2kFXABANm7nYKaj7PakN755hfV0E26z4CJ0FN6GGemnf6E7Y8ZA7VCDjSdnAnqkXf92vdxizZAeApBhsyulvuI2orfMgz9YNqvKRcZB7oPfGB2J0vvZAlmBDEkLPBB2V1nuc8TPmGaLJX1vFNZC5tIzf1YrZC1ox7fXZAbEZBzezZAaFgppf6bXu4YZARP2Bnx5IxWorJgZDZD")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func main() {

	time.LoadLocation("America/Lima")

	ctab := crontab.New()

	err := ctab.AddJob("25 09 * * *", func() {
		fmt.Println("Job")
		// b := []byte("Hola mundo!n")
		// dt := time.Now()
		// formattime := strconv.Itoa(dt.Hour()) +"-"+ strconv.Itoa(dt.Minute()) +"-"+ strconv.Itoa(dt.Second())
		// err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// go getComprobantesElectricos()
	})

	if err != nil {
		log.Fatal(err)
	}

	// ctab.MustAddJob("* * * * *", myFunc)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		getCpeBoletaFactura()

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run("localhost:8080")
}

func task() {
	b := []byte("Hola mundo!\n")
	dt := time.Now()
	// fob := dt.Format("15:04:05")
	formattime := strconv.Itoa(dt.Hour()) + "-" + strconv.Itoa(dt.Minute()) + "-" + strconv.Itoa(dt.Second())
	err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func myFunc() {
	b := []byte("Hola mundo!\n")
	dt := time.Now()
	formattime := strconv.Itoa(dt.Hour()) + "-" + strconv.Itoa(dt.Minute()) + "-" + strconv.Itoa(dt.Second())
	err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
