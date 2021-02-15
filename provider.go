package main

import (
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
)

type Provider struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Entry struct {
	Id int `json:"id"`
	ProviderId int `json:"providerId"`
	Url string `json:"url"`
	Title string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Author string `json:"author"`
	FetchedAt string `json:"fetchedAt"`
}

var feedParser *gofeed.Parser
var fasthttpClient *fasthttp.Client

func InitFeedParser() {
	feedParser = gofeed.NewParser()
	fasthttpClient = &fasthttp.Client{}
}

func ProviderInsert(provider *Provider) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO providers ('name', 'url') VALUES (?, ?)")
	defer stmt.Close()
	if err != nil {
		return -1, err
	}
	op, err := stmt.Exec(provider.Name, provider.Url)
	if err != nil {
		return -1, err
	}
	id, err := op.LastInsertId()
	return id, nil
}

func ProviderGetAll() ([]Provider, error) {
	rows, err := DB.Query("SELECT id, name, url FROM providers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var providers []Provider
	for rows.Next() {
		var provider Provider
		err = rows.Scan(&provider.Id, &provider.Name, &provider.Url)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func ProviderGetById(id int64) (Provider, error) {
	stmt, err := DB.Prepare("SELECT id, name, url FROM providers WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return Provider{}, err
	}
	var provider Provider
	err = stmt.QueryRow(id).Scan(&provider.Id, &provider.Name, &provider.Url)
	if err != nil {
		return Provider{}, err
	}
	return provider, nil
}

func ProviderGetRawEntries(id int64) (*gofeed.Feed, error) {
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
