package models

import (
	"errors"

	"cleanrss/utils"
)


type Provider struct {
	Id       int64  `json:"id" db:"id" validate:"required|isdefault"`
	Name     string `json:"name" db:"name" validate:"required"`
	Url      string `json:"url" db:"url" validate:"required"`
	ParserId int64  `json:"parserId" db:"parser_id" validate:"required|isdefault"`
	IsDeleted bool	`json:"is_deleted" db:"is_deleted"`
}

func ProviderGetAll() ([]Provider, error) {
	providers := []Provider{}
	err := utils.DB.Select(&providers, "SELECT * FROM providers WHERE is_deleted = 0 ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func ProviderGetById(id int64) (Provider, error) {
	provider := Provider{}
	err := utils.DB.Get(&provider, "SELECT * FROM providers WHERE id = $1 AND is_deleted = 0 LIMIT 1", id)
	if err != nil {
		return provider, err
	}
	return provider, nil
}

func ProviderInsert(provider *Provider) (int64, error) {
	if provider == nil {
		return -1, errors.New("Parameter is null")
	}
	res, err := utils.DB.NamedExec("INSERT INTO providers (name, url, parser_id) VALUES (:name, :url, :parser_id)", provider)
	if err != nil {
		return -1, err
	}
	insertedId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return insertedId, nil
}

func ProviderUpdate(provider *Provider) (int64, error) {
	if provider == nil {
		return -1, errors.New("Parameter is null")
	}
	res, err := utils.DB.NamedExec("UPDATE cleaners SET name = :name, url = :url, parser_id = :parser_id WHERE id = :id AND is_deleted = 0", provider)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return affected, nil
}

func ProviderDelete(id int64) (int64, error) {
	res, err := utils.DB.Exec("UPDATE providers SET is_deleted = 1 WHERE id = $1 AND is_deleted = 0", id)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return affected, nil
}