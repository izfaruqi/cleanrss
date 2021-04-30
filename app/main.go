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
	"cleanrss/static"
	"time"

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
	ticker := time.NewTicker(1 * time.Second)
	notificationService, notificationHandler := ws.NewWSNotificationService()

	providerRepository := providerRepo.NewSqliteProviderRepository(db)
	providerUsecase := provider.NewProviderUsecase(providerRepository)
	entryRepository := entryRepo.NewSqliteEntryRepository(db)
	entryUsecase := entry.NewEntryUsecase(entryRepository, entryWebExtRepo.NewWebExtEntryRepository(httpClient, entryRepository, providerUsecase), providerRepository, notificationService)

	mainServer.Mount("/api/provider", providerHttp.NewProviderHTTPHandler(providerUsecase))
	mainServer.Mount("/api/cleaner",
		cleanerHttp.NewCleanerHTTPHandler(cleaner.NewCleanerUsecase(cleanerRepo.NewSqliteCleanerRepository(db), cleanerWebExtRepo.NewWebExtCleanerRepository(httpClient))),
	)
	mainServer.Mount("/api/entry",
		entryHttp.NewEntryHTTPHandler(entryUsecase),
	)
	mainServer.Mount("/api/ws", notificationHandler)
	mainServer.Mount("/", static.NewServeStaticHTTPHandler())

	go func() {
		err := mainServer.Listen("localhost:1337", &wg)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Main server will start on http://localhost:1337")
	tickerEntryUpdater.NewTickerEntryUpdater(ticker, entryUsecase)

	wg.Wait()
}
