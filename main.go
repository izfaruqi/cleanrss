package main

import (
	"log"
)

func main(){
	log.Println("CleanRSS Server starting...")
	DBInit()
	InitFeedParser()
	InitEntryUpdater()

	defer StopEntryUpdater()
	defer DB.Close()

	ServerInit()
}