package models

import (
	"cleanrss/controllers/ws"
	"cleanrss/utils"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
)

type Entry struct {
	Id          int64  `json:"id" db:"id"`
	ProviderId  int64  `json:"providerId" db:"provider_id"`
	Url         string `json:"url" db:"url"`
	Title       string `json:"title" db:"title"`
	PublishedAt int64  `json:"publishedAt" db:"published_at"`
	Author      string `json:"author" db:"author"`
	FetchedAt   int64  `json:"fetchedAt" db:"fetched_at"`
	Json        string `json:"json" db:"json"`
}

func getRawEntriesFromProvider(id int64) (feed *gofeed.Feed, err error) {
	var url string
	err = utils.DB.Get(&url, "SELECT url FROM providers WHERE id = $1 AND is_deleted = 0", id)
	if err != nil {
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	utils.FasthttpClient.Do(req, resp)

	feedParser := gofeed.NewParser()
	feed, err = feedParser.ParseString(strings.TrimSpace(string(resp.Body())))
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func EntryDBRefreshFromProvider(id int64) error {
	feed, err := getRawEntriesFromProvider(id)
	if err != nil {
		return err
	}

	previousEntries := []Entry{}
	err = utils.DB.Select(&previousEntries, "SELECT id, url FROM entries WHERE provider_id=$1 ORDER BY published_at DESC LIMIT $2", id, feed.Len()*2)
	if err != nil {
		return err
	}

	toInsert := make([]Entry, 0, feed.Len())
	timestampNow := time.Now().Unix()

	for _, item := range feed.Items {
		isUpdate := false
		jsonItem, err := json.Marshal(item)
		if err != nil {
			return err
		}
		for _, prev := range previousEntries {
			if prev.Url == item.Link {
				_, err := utils.DB.NamedExec("UPDATE entries SET url = :url, title = :title, published_at = :published_at, author = :author, fetched_at = :fetched_at, json = :json WHERE id = :id", Entry{prev.Id, id, item.Link, item.Title, item.PublishedParsed.Unix(), item.Author.Name, timestampNow, string(jsonItem)})
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
		_, err := utils.DB.NamedExec("INSERT INTO entries (provider_id, url, title, published_at, author, fetched_at, json) VALUES (:provider_id, :url, :title, :published_at, :author, :fetched_at, :json)", toInsert)
		if err != nil {
			return err
		}
	}

	log.Println("Finished updating provider #" + strconv.FormatInt(id, 10))
	ws.WSNotifications <- ws.Notification{Code: "ENTRY_UPDATE_FINISH", Payload: strconv.FormatInt(id, 10)}
	return nil
}

func EntryGetFromDB(providerId int64, limit int, offset int, includeRawJson bool) (*[]Entry, error) {
	entries := []Entry{}
	var err error
	if !includeRawJson {
		if providerId == -1 {
			err = utils.DB.Select(&entries, "SELECT id, provider_id, url, title, published_at, author, fetched_at FROM entries ORDER BY published_at DESC LIMIT $2 OFFSET $3", limit, offset)
		} else {
			err = utils.DB.Select(&entries, "SELECT id, provider_id, url, title, published_at, author, fetched_at FROM entries WHERE provider_id=$1 ORDER BY published_at DESC LIMIT $2 OFFSET $3", providerId, limit, offset)
		}
	} else {
		if providerId == -1 {
			err = utils.DB.Select(&entries, "SELECT * FROM entries ORDER BY published_at DESC LIMIT $2 OFFSET $3", limit, offset)
		} else {
			err = utils.DB.Select(&entries, "SELECT * FROM entries WHERE provider_id=$1 ORDER BY published_at DESC LIMIT $2 OFFSET $3", providerId, limit, offset)
		}
	}
	if err != nil {
		return nil, err
	}
	return &entries, nil
}

func EntrySearch(query string, providerId int64) (*[]Entry, error) {
	entries := []Entry{}
	var err error
	if providerId != -1 {
		err = utils.DB.Select(&entries, "SELECT id, provider_id, url, title, published_at, author, fetched_at FROM entries WHERE (title LIKE $1) AND provider_id = $2", "%"+query+"%", providerId)
	} else {
		err = utils.DB.Select(&entries, "SELECT id, provider_id, url, title, published_at, author, fetched_at FROM entries WHERE (title LIKE $1)", "%"+query+"%")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &entries, nil
}
