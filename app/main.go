package main

import (
	tickerEntryUpdater "cleanrss/entry/delivery/ticker"
	"cleanrss/infrastructure"
	"time"

	"cleanrss/cleaner"
	cleanerHttp "cleanrss/cleaner/delivery/http"
	cleanerRepo "cleanrss/cleaner/repository/sqlite"
	cleanerWebExtRepo "cleanrss/cleaner/repository/web_ext"
	"cleanrss/entry"
	entryHttp "cleanrss/entry/delivery/http"
	entryRepo "cleanrss/entry/repository/sqlite"
	entryWebExtRepo "cleanrss/entry/repository/web_ext"
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

	httpClient := infrastructure.NewHTTPClient()
	mainServer := infrastructure.NewHTTPServer()
	proxyServer := infrastructure.NewHTTPServer()
	ticker := time.NewTicker(1 * time.Second)

	providerRepository := providerRepo.NewSqliteProviderRepository(db)
	providerUsecase := provider.NewProviderUsecase(providerRepository)
	entryRepository := entryRepo.NewSqliteEntryRepository(db)
	entryUsecase := entry.NewEntryUsecase(entryRepository, entryWebExtRepo.NewWebExtEntryRepository(httpClient, entryRepository, providerUsecase), providerRepository)
	providerHttp.NewProviderHttpHandler(mainServer.Group("/api/provider"), providerUsecase)
	cleanerHttp.NewCleanerHttpHandler(mainServer.Group("/api/cleaner"), cleaner.NewCleanerUsecase(cleanerRepo.NewSqliteCleanerRepository(db), cleanerWebExtRepo.NewWebExtCleanerRepository(httpClient)))
	entryHttp.NewEntryHttpRouter(mainServer.Group("/api/entry"), entryUsecase)

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
	tickerEntryUpdater.NewTickerEntryUpdater(ticker, entryUsecase)

	wg.Wait()
}
