package main

import (
	"log"
)

func main(){
	log.Println("CleanRSS Server starting...")
	DBInit()
	InitFeedParser()
	ServerInit()

	defer DB.Close()
}