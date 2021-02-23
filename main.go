package main

import (
	"log"

	"cleanrss/services"
	"cleanrss/utils"
)

func main(){
	log.Println("CleanRSS Server starting...")
	utils.DBInit()
	utils.HttpClientInit()
	services.EntryUpdaterInit()

	defer utils.DB.Close()
	defer log.Println("CleanRSS Server shutting down...")
	
	ServerInit()
}