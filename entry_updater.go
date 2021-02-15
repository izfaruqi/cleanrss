package main

import (
	"log"

	"github.com/robfig/cron"
)

var EntryUpdaterCron *cron.Cron

func InitEntryUpdater() {
	EntryUpdaterCron = cron.New()
	//EntryUpdaterCron.AddFunc("0 * * * * *", func() { ProviderRefreshEntriesForDB(*put id here*) })
	EntryUpdaterCron.Start()
	log.Println("Started Entry Updater")
}

func StopEntryUpdater(){
	EntryUpdaterCron.Stop()
}