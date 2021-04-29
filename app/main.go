package main

import (
	"cleanrss/cleaner"
	cleanerHttp "cleanrss/cleaner/delivery/http"
	cleanerRepo "cleanrss/cleaner/repository/sqlite"
	cleanerWebExtRepo "cleanrss/cleaner/repository/web_ext"
	"cleanrss/entry"
	entryHttp "cleanrss/entry/delivery/http"
	tickerEntryUpdater "cleanrss/entry/delivery/ticker"
	entryRepo "cleanrss/entry/repository/sqlite"
	entryWebExtRepo "cleanrss/entry/repository/web_ext"
	"cleanrss/infrastructure"
	"cleanrss/infrastructure/notification/ws"
	"cleanrss/provider"
	providerHttp "cleanrss/provider/delivery/http"
	providerRepo "cleanrss/provider/repository/sqlite"
	"time"

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
	mainChiServer := infrastructure.NewHTTPChiServer()
	proxyServer := infrastructure.NewHTTPServer()
	ticker := time.NewTicker(1 * time.Second)
	notificationService, notificationHandler := ws.NewWSNotificationService()

	providerRepository := providerRepo.NewSqliteProviderRepository(db)
	providerUsecase := provider.NewProviderUsecase(providerRepository)
	entryRepository := entryRepo.NewSqliteEntryRepository(db)
	entryUsecase := entry.NewEntryUsecase(entryRepository, entryWebExtRepo.NewWebExtEntryRepository(httpClient, entryRepository, providerUsecase), providerRepository, notificationService)
	providerHttp.NewProviderHttpHandler(mainServer.Group("/api/provider"), providerUsecase)
	cleanerHttp.NewCleanerHttpHandler(mainServer.Group("/api/cleaner"), cleaner.NewCleanerUsecase(cleanerRepo.NewSqliteCleanerRepository(db), cleanerWebExtRepo.NewWebExtCleanerRepository(httpClient)))
	entryHttp.NewEntryHttpRouter(mainServer.Group("/api/entry"), entryUsecase)

	mainChiServer.Mount("/api/provider", providerHttp.NewProviderHTTPChiHandler(providerUsecase))
	mainChiServer.Mount("/api/cleaner",
		cleanerHttp.NewCleanerHTTPChiHandler(cleaner.NewCleanerUsecase(cleanerRepo.NewSqliteCleanerRepository(db), cleanerWebExtRepo.NewWebExtCleanerRepository(httpClient))),
	)
	mainChiServer.Mount("/api/entry",
		entryHttp.NewEntryHTTPChiHandler(entryUsecase),
	)
	mainChiServer.Mount("/api/ws", notificationHandler)

	proxyHttp.NewProxyHandler(proxyServer.App, httpClient, "/proxy", "http://localhost:1338")

	go func() {
		err := mainChiServer.Listen("localhost:1336", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Chi server will start on http://localhost:1336")
	go func() {
		err := mainServer.Listen("localhost:1337", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Main server will start on http://localhost:1337")
	go func() {
		err := proxyServer.Listen("localhost:1338", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Proxy server will start on http://localhost:1338")
	tickerEntryUpdater.NewTickerEntryUpdater(ticker, entryUsecase)

	wg.Wait()
}
