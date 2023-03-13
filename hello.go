package main

import (
	// "github.com/go-co-op/gocron"
	"github.com/mileusna/crontab"
	"log"
	"io/ioutil"
	// "fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
)
 
func main(){
	time.LoadLocation("America/Lima")	
	
    ctab := crontab.New()

    err := ctab.AddJob("32 22 * * *", func(){
		b := []byte("Hola mundo!\n")
		dt := time.Now()
		formattime := strconv.Itoa(dt.Hour()) +"-"+ strconv.Itoa(dt.Minute()) +"-"+ strconv.Itoa(dt.Second()) 
		err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
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

func task(){
		b := []byte("Hola mundo!\n")
		dt := time.Now()
		// fob := dt.Format("15:04:05")
		formattime := strconv.Itoa(dt.Hour()) +"-"+ strconv.Itoa(dt.Minute()) +"-"+ strconv.Itoa(dt.Second()) 
		err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
}

func myFunc() {
	b := []byte("Hola mundo!\n")
	dt := time.Now()
	formattime := strconv.Itoa(dt.Hour()) +"-"+ strconv.Itoa(dt.Minute()) +"-"+ strconv.Itoa(dt.Second()) 
	err := ioutil.WriteFile(formattime+"-personal.txt", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}