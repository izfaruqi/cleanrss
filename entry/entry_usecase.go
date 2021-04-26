package entry

import (
	"cleanrss/domain"
	"log"
	"strconv"
)

type entryUsecase struct {
	r  domain.EntryRepository
	we domain.WebExtEntryRepository
	rp domain.ProviderRepository
}

func NewEntryUsecase(r domain.EntryRepository, we domain.WebExtEntryRepository, rp domain.ProviderRepository) domain.EntryUsecase {
	return entryUsecase{r: r, we: we, rp: rp}
}

func (e entryUsecase) GetById(id int64, withJson bool) (domain.Entry, error) {
	panic("implement me")
}

func (e entryUsecase) GetAll(withJson bool) (*[]domain.Entry, error) {
	panic("implement me")
}

func (e entryUsecase) GetByQuery(query string, dateFrom int64, dateUntil int64, providerId int64, limit int64, offset int64, withJson bool) ([]domain.Entry, error) {
	return e.r.GetByQuery(query, dateFrom, dateUntil, providerId, limit, offset, withJson, false)
}

func (e entryUsecase) TriggerRefresh(providerId int64) error {
	entries, entriesLen, err := e.we.GetRawEntriesByProviderId(providerId)
	if err != nil {
		return err
	}

	previousEntries, err := e.r.GetByQuery("", -1, -1, providerId, int64(entriesLen*2), 0, false, true)
	if err != nil {
		return err
	}

	toInsert := make([]domain.Entry, 0, entriesLen)

	for _, item := range entries {
		isUpdate := false
		for _, prev := range previousEntries {
			if prev.Url == item.Url {
				item.Id = prev.Id
				err = e.r.Update(item)
				if err != nil {
					return err
				}
				isUpdate = true
				break
			}
		}
		if !isUpdate {
			toInsert = append(toInsert, item)
		}
	}
	if len(toInsert) > 0 {
		err := e.r.BulkInsert(toInsert)
		if err != nil {
			return err
		}
	}

	log.Println("Finished updating provider #" + strconv.FormatInt(providerId, 10))
	//ws.WSNotifications <- ws.Notification{Code: "ENTRY_UPDATE_FINISH", Payload: strconv.FormatInt(providerId, 10)}
	return nil
}

func (e entryUsecase) TriggerRefreshAll() error {
	providers, err := e.rp.GetAll()
	if err != nil {
		return err
	}
	for _, provider := range *providers {
		provider := provider
		go func() { // TODO: Have a notification system for errors
			err := e.TriggerRefresh(provider.Id)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	return nil
}
