package domain

import "encoding/json"

type Cleaner struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Rules     *CleanerRules
	RulesJson string `json:"rulesJson" db:"rulesJson"`
	IsDeleted bool   `json:"is_deleted" db:"isDeleted"`
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
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider Provider) error
	Update(id int64, provider Provider) error
	Delete(id int64) error
}

type CleanerRepository interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider Provider) error
	Update(id int64, provider Provider) error
	Delete(id int64) error
}
