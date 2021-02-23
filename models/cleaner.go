package models

import (
	"bytes"
	"cleanrss/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

type Cleaner struct {
	Id        int64  `json:"id" db:"id" validate:"required|isdefault"`
	Name      string `json:"name" db:"name" validate:"required"`
	RulesJson string  `json:"rulesJson" db:"rules_json" validate:"required"`
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}

func CleanerGetAll() ([]Cleaner, error) {
	cleaners := []Cleaner{}
	err := utils.DB.Select(&cleaners, "SELECT * FROM cleaners WHERE is_deleted = 0 ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return cleaners, nil
}

func CleanerGetById(id int64) (Cleaner, error) {
	cleaner := Cleaner{}
	err := utils.DB.Get(&cleaner, "SELECT * FROM cleaners WHERE id = $1 AND is_deleted = 0 LIMIT 1", id)
	if err != nil {
		return cleaner, err
	}
	return cleaner, nil
}

func CleanerInsert(cleaner *Cleaner) (int64, error) {
	if cleaner == nil {
		return -1, errors.New("Parameter is null")
	}
	res, err := utils.DB.NamedExec("INSERT INTO cleaners (name, rules_json) VALUES (:name, :rules_json)", cleaner)
	if err != nil {
		return -1, err
	}
	insertedId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return insertedId, nil
}

func CleanerUpdate(cleaner *Cleaner) (int64, error) {
	if cleaner == nil {
		return -1, errors.New("Parameter is null")
	}
	res, err := utils.DB.NamedExec("UPDATE cleaners SET name = :name, rules_json = :rules_json WHERE id = :id AND is_deleted = 0", cleaner)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return affected, nil
}

func CleanerDelete(id int64) (int64, error) {
	res, err := utils.DB.Exec("UPDATE cleaners SET is_deleted = 1 WHERE id = $1 AND is_deleted = 0", id)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return affected, nil
}


func getRawPage(url string, ua string) *bytes.Reader {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(url)
	req.Header.Set("User-Agent", ua)

	utils.FasthttpClient.DoRedirects(req, res, 20)
	return bytes.NewReader(res.Body())
}

func cleanPage(url string, parserJson map[string]interface{}) (string, error) {
	var requestRules, htmlRules map[string]interface{}
	if parserJson["request"] == nil && parserJson["html"] == nil {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, getRawPage(url, "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 Mobile Safari/537.36"))
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	requestRules = parserJson["request"].(map[string]interface{})
	htmlRules = parserJson["html"].(map[string]interface{})

	var pageBodyReader *bytes.Reader

	if requestRules["mobileUA"] != nil {
		if requestRules["mobileUA"].(bool) {
			pageBodyReader = getRawPage(url, "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 Mobile Safari/537.36")
		}
	} else {
		pageBodyReader = getRawPage(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	}
	
	doc, err := goquery.NewDocumentFromReader(pageBodyReader)
	if err != nil {
		return "", err
	}

	var rootNode *goquery.Selection
	rootRules := htmlRules["root"].([]interface{})
	for _, rootRule := range rootRules {
		rootNode = doc.Find(rootRule.(string)).First()
		if rootNode.Length() != 0 {
			break
		}
	}
	
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
	return outStr, nil
}

func CleanerGetPage(entryId int64) (string, error) {
	rows, err := utils.DB.Queryx("SELECT entries.url, cleaners.rules_json FROM entries LEFT JOIN providers ON entries.provider_id = providers.id LEFT JOIN cleaners ON providers.parser_id = cleaners.id WHERE entries.id = $1 LIMIT 1", entryId)
	if err != nil {
		log.Fatalln(err)
	}
	rows.Next()
	cols, _ := rows.SliceScan()
	url := cols[0].(string)
	var parserJson interface{}
	err = json.Unmarshal([]byte(cols[1].(string)), &parserJson)
	cleanedPage, err := cleanPage(url, parserJson.(map[string]interface{}))
	if err != nil {
		return "", err
	}
	return cleanedPage, nil
}