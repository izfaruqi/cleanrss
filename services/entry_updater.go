package services

import (
	"cleanrss/models"
	"cleanrss/utils"
	"log"

	"github.com/robfig/cron"
)

var entryUpdateTimer *cron.Cron

// TODO: Implement lock system to prevent potential race conditions when a manual update is trigger when update timer is running.

func EntryUpdaterInit(){
	if entryUpdateTimer != nil {
		entryUpdateTimer.Stop()
	}
	entryUpdateTimer = cron.New()
	entryUpdateTimer.AddFunc("0 30 * * * *", refreshEntriesFromProviders)
	entryUpdateTimer.Start()
}

func refreshEntriesFromProviders() {
	log.Println("Updating entries...")
	providerIds := refreshEntryUpdaterProviders()
	for _, providerId := range providerIds {
		go launchEntryDBRefresh(providerId)
	}
}

func launchEntryDBRefresh(providerId int64){
	err := models.EntryDBRefreshFromProvider(providerId)
	if err != nil {
		log.Println(err.Error())
	}
}


func refreshEntryUpdaterProviders() []int64 {
	providerIds := make([]int64, 0)
	err := utils.DB.Select(&providerIds, "SELECT id FROM providers WHERE is_deleted = 0")
	if err != nil {
		log.Fatalln(err)
	}
	return providerIds
}