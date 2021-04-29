package web_ext

import (
	"cleanrss/domain"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type webExtEntryRepository struct {
	client *http.Client
	re     domain.EntryRepository
	rp     domain.ProviderRepository
}

func (w webExtEntryRepository) GetRawEntriesByProviderId(providerId int64) ([]domain.Entry, int, error) {
	provider, err := w.rp.GetById(providerId)
	if err != nil {
		return nil, 0, err
	}
	url := provider.Url
	req, err := http.NewRequest("GET", url, nil)
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	timestampNow := time.Now().Unix()
	feedParser := gofeed.NewParser()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	feed, err := feedParser.ParseString(strings.TrimSpace(string(body))) // Some sites has feeds that contains some kind of trailing spaces.
	if err != nil {
		return nil, 0, err
	}
	entries := make([]domain.Entry, feed.Len())
	for i, item := range feed.Items {
		jsonItem, err := json.Marshal(item)
		if err != nil {
			return nil, 0, err
		}
		entries[i] = domain.Entry{Id: -1, ProviderId: providerId, Url: item.Link, Title: item.Title, PublishedAt: item.PublishedParsed.Unix(), Author: item.Author.Name, FetchedAt: timestampNow, Json: string(jsonItem)}
	}
	return entries, feed.Len(), nil
}

func NewWebExtEntryRepository(httpClient *http.Client, re domain.EntryRepository, rp domain.ProviderRepository) domain.WebExtEntryRepository {
	return webExtEntryRepository{client: httpClient, re: re, rp: rp}
}
