package sqlite

import (
	"cleanrss/domain"
	"errors"
	"github.com/jmoiron/sqlx"
)

type sqliteCleanerRepository struct {
	DB *sqlx.DB
}

func NewSqliteCleanerRepository(db *sqlx.DB) domain.CleanerRepository {
	return sqliteCleanerRepository{
		DB: db,
	}
}

func (m sqliteCleanerRepository) GetById(id int64) (domain.Cleaner, error) {
	cleaner := domain.Cleaner{}
	err := m.DB.Get(&cleaner, "SELECT * FROM cleaners WHERE id = $1 AND is_deleted = 0 LIMIT 1", id)
	if err != nil {
		return cleaner, err
	}
	return cleaner, nil
}

func (m sqliteCleanerRepository) GetAll() (*[]domain.Cleaner, error) {
	var cleaners []domain.Cleaner
	err := m.DB.Select(&cleaners, "SELECT * FROM cleaners WHERE is_deleted = 0 ORDER BY id")
	if err != nil {
		return nil, err
	}
	return &cleaners, nil
}

func (m sqliteCleanerRepository) Insert(cleaner *domain.Cleaner) error {
	if cleaner == nil {
		return errors.New("Parameter is null")
	}
	res, err := m.DB.NamedExec("INSERT INTO cleaners (name, rules_json) VALUES (:name, :rules_json)", cleaner)
	if err != nil {
		return err
	}
	insertedId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cleaner.Id = insertedId
	return nil
}

func (m sqliteCleanerRepository) Update(cleaner domain.Cleaner) error {
	res, err := m.DB.NamedExec("UPDATE cleaners SET name = :name, rules_json = :rules_json WHERE id = :id AND is_deleted = 0", cleaner)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("id not found")
	}
	return nil
}

func (m sqliteCleanerRepository) Delete(id int64) error {
	res, err := m.DB.Exec("UPDATE cleaners SET is_deleted = 1 WHERE id = $1 AND is_deleted = 0", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("id not found")
	}
	return nil
}
