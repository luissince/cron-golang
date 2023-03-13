package main

import (
	// "github.com/go-co-op/gocron"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/crontab"
)

func getComprobantesElectricos() {
	res, err := http.Get("http://localhost:9000/app/controller/VentaController.php?type=listComprobantes&opcion=0&busqueda=&fechaInicial=&fechaFinal=&comprobante=0&estado=0&posicionPagina=0&filasPorPagina=10")
	if err != nil {
		fmt.Println("Ocurrio un error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ocurrio un error")
	}

	fmt.Printf("%s", body)

}

func main() {

	time.LoadLocation("America/Lima")

	ctab := crontab.New()

	err := ctab.AddJob("25 09 * * *", func() {
		fmt.Println("Job")
		// b := []byte("Hola mundo!\n")
		// dt := time.Now()
		// formattime := strconv.Itoa(dt.Hour()) +"-"+ strconv.Itoa(dt.Minute()) +"-"+ strconv.Itoa(dt.Second())
		// err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		go getComprobantesElectricos()
	})

	if err != nil {
		log.Fatal(err)
	}

	ctab.MustAddJob("* * * * *", myFunc)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
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
