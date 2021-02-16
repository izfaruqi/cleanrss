package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)


func GetCleanPage(entryId int64) string {
	rows, err := DB.Queryx("SELECT entries.url, parsers.rules_json FROM entries LEFT JOIN providers ON entries.provider_id = providers.id LEFT JOIN parsers ON providers.parser_id = parsers.id WHERE entries.id = $1 LIMIT 1", entryId)
	if err != nil {
		log.Fatalln(err)
	}
	rows.Next()
	cols, _ := rows.SliceScan()
	url := cols[0].(string)
	var parserJson interface{}
	err = json.Unmarshal([]byte(cols[1].(string)), &parserJson)
	cleanedPage := parsePage(url, parserJson.(map[string]interface{}))
	return cleanedPage
}

func parsePage(url string, parserJson map[string]interface{}) string {
	requestRules := parserJson["request"].(map[string]interface{})
	htmlRules := parserJson["html"].(map[string]interface{})

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(url)
	if requestRules["mobileUA"] != nil {
		if requestRules["mobileUA"].(bool) {
			req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 Mobile Safari/537.36")
		}
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	}
	fasthttpClient.DoRedirects(req, res, 20)
	
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {

	}
	
	rootNode := doc.Find(htmlRules["root"].(string)).First()

	if htmlRules["noscript"] != nil {
		if htmlRules["noscript"].(bool) {
			doc.Find("noscript").Each(func(i int, s *goquery.Selection){
				s.Parent().AppendHtml(s.Text())
			})
		}
	}

	if htmlRules["remove"] != nil {
		for _, toRemove := range htmlRules["remove"].([]interface{}) {
			doc.Find(toRemove.(string)).Each(func(i int, s *goquery.Selection){
				s.Remove()
		 })
		}
	}
	
	outStr, _ := goquery.OuterHtml(rootNode)
	return outStr
}