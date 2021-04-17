package main

import (
	"log"
	"sync"

	"cleanrss/services"
	"cleanrss/utils"
)

func main() {
	log.Println("CleanRSS Server starting...")
	utils.DBInit()
	utils.HttpClientInit()
	services.EntryUpdaterInit()

	defer utils.DB.Close()
	defer log.Println("CleanRSS Server shutting down...")

	var wg sync.WaitGroup
	wg.Add(2)

	go ServerInit(&wg)
	go ProxyServerInit("localhost:3333", &wg)

	wg.Wait()
}
