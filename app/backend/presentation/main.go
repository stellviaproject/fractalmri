package main

import (
	"fmt"
	app "fractalmri/app/backend/application"
	config "fractalmri/app/backend/domain/config"
	"log"
	"os"
	"strconv"
)

func main() {
	//Obtener el puerto en el que escuchar desde la variable de entorno PORT
	portStr := os.Getenv("PORT")
	var err error
	var port int
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Println(err)
		}
	}
	var configuration *config.Configuration
	configuration, err = config.LoadConfig("./config.json")
	if err != nil {
		log.Println(err)
		configuration = config.Default()
		if err := configuration.SaveConfig(); err != nil {
			log.Fatalln(err)
		}
	}
	if portStr == "" {
		port = configuration.Port
	}
	app := app.NewApp()
	app.SetConfig(configuration)
	app.Run(fmt.Sprintf(":%d", port))
}
