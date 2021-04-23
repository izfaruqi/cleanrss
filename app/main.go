package main

import (
	"cleanrss/infrastructure"

	"cleanrss/cleaner"
	cleanerHttp "cleanrss/cleaner/delivery/http"
	cleanerRepo "cleanrss/cleaner/repository/sqlite"
	"cleanrss/provider"
	providerHttp "cleanrss/provider/delivery/http"
	providerRepo "cleanrss/provider/repository/sqlite"

	proxyHttp "cleanrss/proxy/delivery/http"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	db, err := infrastructure.NewDB()
	if err != nil {
		log.Println(err)
		return
	}
	mainServer := infrastructure.NewHTTPServer()
	proxyServer := infrastructure.NewHTTPServer()

	providerHttp.NewProviderHttpHandler(mainServer.Group("/provider"), provider.NewProviderUsecase(providerRepo.NewSqliteProviderRepository(db)))
	cleanerHttp.NewCleanerHttpHandler(mainServer.Group("/cleaner"), cleaner.NewCleanerUsecase(cleanerRepo.NewSqliteCleanerRepository(db)))

	proxyHttp.NewProxyHandler(proxyServer.App, "/proxy", "http://localhost:8081")

	go func() {
		err := mainServer.Listen("localhost:8080", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Main server will start on http://localhost:8080")
	go func() {
		err := proxyServer.Listen("localhost:8081", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Proxy server will start on http://localhost:8081")

	wg.Wait()
}
