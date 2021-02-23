package main

import (
	"log"

	"cleanrss/utils"
)

func main(){
	log.Println("CleanRSS Server starting...")
	utils.DBInit()
	utils.HttpClientInit()
	
	defer utils.DB.Close()

	ServerInit()
}