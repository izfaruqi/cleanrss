package main

import (
	"cleanrss/cleaner"
	cleanerHttp "cleanrss/cleaner/delivery/http"
	cleanerRepo "cleanrss/cleaner/repository/sqlite"
	cleanerWebExtRepo "cleanrss/cleaner/repository/web_ext"
	"cleanrss/entry"
	entryHttp "cleanrss/entry/delivery/http"
	entryRepo "cleanrss/entry/repository/sqlite"
	entryWebExtRepo "cleanrss/entry/repository/web_ext"
	"cleanrss/infrastructure"
	"cleanrss/infrastructure/notification/ws"
	"cleanrss/provider"
	providerHttp "cleanrss/provider/delivery/http"
	providerRepo "cleanrss/provider/repository/sqlite"
	"cleanrss/static"
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	db, err := infrastructure.NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	httpClient := infrastructure.NewHTTPClient()
	mainServer := infrastructure.NewHTTPServer()
	mainServerShutdownCtx, mainServerShutdown := context.WithCancel(context.Background())
	cronScheduler := cron.New()
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

	var waitSystray sync.WaitGroup
	infoSystray := infrastructure.SystrayInfo{}
	waitSystray.Add(1)
	go infrastructure.RunSystray(&infoSystray, &waitSystray)
	go func() {
		waitSystray.Wait()
		infoSystray.SetListeningAddress("localhost:1337")
		infoSystray.SetOnQuitClicked(func() {
			mainServerShutdown()
		})
		err := mainServer.Listen("localhost:1337", &wg, mainServerShutdownCtx)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Main server will start on http://localhost:1337")

	_, err = cronScheduler.AddFunc("0 * * * *", func() {
		log.Println("Updating all providers...")
		err := entryUsecase.TriggerRefreshAll()
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
	cronScheduler.Start()

	wg.Wait()
}
