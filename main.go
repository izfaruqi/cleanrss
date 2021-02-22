package main

import (
	"log"

	"cleanrss/utils"
)

func main(){
	log.Println("CleanRSS Server starting...")
	utils.DBInit()
	//InitFeedParser()
	//InitEntryUpdater()

	//defer StopEntryUpdater()
	defer utils.DB.Close()

	ServerInit()
}