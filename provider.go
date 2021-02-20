package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
)

type Provider struct {
	Id   int64    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Url  string `json:"url" db:"url"`
	ParserId int64 `json:"parserId" db:"parser_id"`
}

type Entry struct {
	Id int64 `json:"id" db:"id"`
	ProviderId int64 `json:"providerId" db:"provider_id"`
	Url string `json:"url" db:"url"`
	Title string `json:"title" db:"title"`
	PublishedAt int64 `json:"publishedAt" db:"published_at"`
	Author string `json:"author" db:"author"`
	FetchedAt int64 `json:"fetchedAt" db:"fetched_at"`
	Json string `json:"json" db:"json"`
}

var feedParser *gofeed.Parser
var fasthttpClient *fasthttp.Client

func InitFeedParser() {
	feedParser = gofeed.NewParser()
	fasthttpClient = &fasthttp.Client{}
}

func ProviderInsert(provider *Provider) (int64, error) {
	res, err := DB.NamedExec("INSERT INTO providers ('name', 'url', 'parser_id') VALUES (:name, :url, :parser_id)", provider)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}
	return id, nil
}

func ProviderUpdate(provider *Provider) (error) {
	res, err := DB.NamedExec("UPDATE providers SET name = :name, url = :url, parser_id = :parser_id WHERE id = :id", provider)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func ProviderGetAll() ([]Provider, error) {
	providers := []Provider{}
	err := DB.Select(&providers, "SELECT * FROM providers")
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func ProviderGetById(id int64) (Provider, error) {
	provider := Provider{}
	err := DB.Get(&provider, "SELECT * FROM providers WHERE id = $1", id)
	if err != nil {
		return provider, err
	}
	return provider, nil
}

func ProviderGetFreshEntries(id int64) (*gofeed.Feed, error) {
	stmt, err := DB.Prepare("SELECT url FROM providers WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var url string
	err = stmt.QueryRow(id).Scan(&url)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	fasthttpClient.Do(req, resp)

	feed, err := feedParser.ParseString(strings.TrimSpace(string(resp.Body())))
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func ProviderRefreshEntriesForDB(id int64) error {
	fmt.Println("Updating provider #" + strconv.FormatInt(id, 10))

	rawEntries, err := ProviderGetFreshEntries(id)
	if err != nil {
		return err
	}

	previousEntries := []Entry{}
	err = DB.Select(&previousEntries, "SELECT id, url FROM entries WHERE provider_id=$1 ORDER BY published_at DESC LIMIT $2", id, rawEntries.Len())
	if err != nil {
		return err
	}

	toInsert := make([]Entry, 0, rawEntries.Len())
	timestampNow := time.Now().Unix()

	for _, item := range rawEntries.Items {
		isUpdate := false
		jsonItem, err := json.Marshal(item)
		if err != nil {
			return err
		}
		for _, prev := range previousEntries {
			if prev.Url == item.Link {
				_, err := DB.NamedExec("UPDATE entries SET url = :url, title = :title, published_at = :published_at, author = :author, fetched_at = :fetched_at, json = :json WHERE id = :id", Entry{prev.Id, id, item.Link, item.Title, item.PublishedParsed.Unix(), item.Author.Name, timestampNow, string(jsonItem)})
				if err != nil {
					return err
				}
				isUpdate = true
				break
			}
		}
		if !isUpdate {
			toInsert = append(toInsert, Entry{-1, id, item.Link, item.Title, item.PublishedParsed.Unix(), item.Author.Name, timestampNow, string(jsonItem)})
		}
	}
	if len(toInsert) > 0 {
		_, err := DB.NamedExec("INSERT INTO entries (provider_id, url, title, published_at, author, fetched_at, json) VALUES (:provider_id, :url, :title, :published_at, :author, :fetched_at, :json)", toInsert)
		if err != nil {
			return err
		}
	}

	fmt.Println("Finished updating provider #" + strconv.FormatInt(id, 10))
	return nil
}

func ProviderGetDBEntries(providerId int64, limit int) (*[]Entry, error){
	entries := []Entry{}
	err := DB.Select(&entries, "SELECT * FROM entries WHERE provider_id=$1 ORDER BY published_at DESC LIMIT $2", providerId, limit)
	if err != nil {
		return nil, err
	}
	return &entries, nil
}