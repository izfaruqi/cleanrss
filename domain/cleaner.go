package domain

import (
	"encoding/json"
	"io"
)

type Cleaner struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Rules     *CleanerRules
	RulesJson string `json:"rulesJson" db:"rules_json"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
}

type CleanerRules struct {
}

func (c Cleaner) ParseRules() error {
	if c.Rules != nil {
		return nil
	} // Dont parse again if already parsed.
	err := json.Unmarshal([]byte(c.RulesJson), c.Rules)
	if err != nil {
		return err
	}
	return nil
}

func (c Cleaner) RulesToJson() error {
	rulesJson, err := json.Marshal(c.Rules)
	if err != nil {
		return err
	}
	c.RulesJson = string(rulesJson)
	return nil
}

type CleanerUsecase interface {
	GetById(id int64) (Cleaner, error)
	GetAll() (*[]Cleaner, error)
	Insert(cleaner *Cleaner) error
	Update(cleaner Cleaner) error
	Delete(id int64) error
	GetCleanedEntry(entryId int64) (string, error)
}

type CleanerRepository interface {
	GetById(id int64) (Cleaner, error)
	GetAll() (*[]Cleaner, error)
	Insert(cleaner *Cleaner) error
	Update(cleaner Cleaner) error
	Delete(id int64) error
	GetEntryUrlAndCleaner(id int64) (url string, cleaner string, err error)
}

type WebExtCleanerRepository interface {
	GetRawPage(url string, mobileUA bool) (io.Reader, error)
}
