package cleaner

import (
	"cleanrss/domain"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

type cleanerUsecase struct {
	r  domain.CleanerRepository
	we domain.WebExtCleanerRepository
}

func NewCleanerUsecase(r domain.CleanerRepository, we domain.WebExtCleanerRepository) domain.CleanerUsecase {
	return cleanerUsecase{r: r, we: we}
}

func (c cleanerUsecase) GetCleanedEntry(entryId int64) (string, error) {
	url, cleanerString, err := c.r.GetEntryUrlAndCleaner(entryId)
	if err != nil {
		return "", err
	}
	var parserJson interface{}
	err = json.Unmarshal([]byte(cleanerString), &parserJson)
	cleanedPage, err := c.cleanPage(url, parserJson.(map[string]interface{}))
	if err != nil {
		return "", err
	}
	return cleanedPage, nil
}

func (c cleanerUsecase) cleanPage(url string, parserJson map[string]interface{}) (string, error) {
	var requestRules, htmlRules map[string]interface{}
	var err error

	if parserJson["request"] == nil && parserJson["html"] == nil {
		buf := new(strings.Builder)
		rawPage, err := c.we.GetRawPage(url, false)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(buf, rawPage)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	requestRules = parserJson["request"].(map[string]interface{})
	htmlRules = parserJson["html"].(map[string]interface{})

	var pageBodyReader io.Reader

	if requestRules["mobileUA"] != nil {
		pageBodyReader, err = c.we.GetRawPage(url, requestRules["mobileUA"].(bool))
	} else {
		pageBodyReader, err = c.we.GetRawPage(url, false)
	}
	if err != nil {
		return "", nil
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
			rootNode.Find("noscript").Each(func(i int, s *goquery.Selection) {
				s.Parent().AppendHtml(s.Text())
			})
		}
	}

	if htmlRules["remove"] != nil {
		for _, toRemove := range htmlRules["remove"].([]interface{}) {
			rootNode.Find(toRemove.(string)).Each(func(i int, s *goquery.Selection) {
				s.Remove()
			})
		}
	}

	rootNode.Find("a").SetAttr("target", "_blank")

	outStr, err := goquery.OuterHtml(rootNode)
	if err != nil {
		return "", nil
	}
	return outStr, nil
}

func (c cleanerUsecase) GetById(id int64) (domain.Cleaner, error) {
	return c.r.GetById(id)
}

func (c cleanerUsecase) GetAll() (*[]domain.Cleaner, error) {
	return c.r.GetAll()
}

func (c cleanerUsecase) Insert(cleaner *domain.Cleaner) error {
	return c.r.Insert(cleaner)
}

func (c cleanerUsecase) Update(cleaner domain.Cleaner) error {
	return c.r.Update(cleaner)
}

func (c cleanerUsecase) Delete(id int64) error {
	return c.r.Delete(id)
}
