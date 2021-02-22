package models

import (
	"cleanrss/utils"
	"errors"
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