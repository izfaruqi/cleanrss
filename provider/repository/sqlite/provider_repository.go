package sqlite

import (
	"cleanrss/domain"
	"github.com/jmoiron/sqlx"
)

type sqliteProviderRepository struct {
	DB *sqlx.DB
}

func NewSqliteProviderRepository(db *sqlx.DB) domain.ProviderRepository {
	return sqliteProviderRepository{
		DB: db,
	}
}

func (m sqliteProviderRepository) GetById(id int64) (domain.Provider, error) {
	provider := domain.Provider{}
	err := m.DB.Get(&provider, "SELECT * FROM providers WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return provider, err
	}
	return provider, nil
}

func (m sqliteProviderRepository) GetAll() (*[]domain.Provider, error) {
	var providers []domain.Provider
	err := m.DB.Select(&providers, "SELECT * FROM providers WHERE is_deleted = 0 ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return &providers, nil
}

func (m sqliteProviderRepository) Insert(provider *domain.Provider) error {
	res, err := m.DB.NamedExec("INSERT INTO providers (name, url, parser_id) VALUES (:name, :url, :parser_id)", provider)
	if err != nil {
		return err
	}
	insertedId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	provider.Id = insertedId
	return nil
}

func (m sqliteProviderRepository) Update(provider domain.Provider) error {
	res, err := m.DB.NamedExec("UPDATE providers SET name = :name, url = :url, parser_id = :parser_id WHERE id = :id AND is_deleted = 0", provider)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (m sqliteProviderRepository) Delete(id int64) error {
	res, err := m.DB.Exec("UPDATE providers SET is_deleted = 1 WHERE id = $1 AND is_deleted = 0", id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
