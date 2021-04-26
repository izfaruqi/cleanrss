package sqlite

import (
	"cleanrss/domain"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type sqliteEntryRepository struct {
	db *sqlx.DB
}

func NewSqliteEntryRepository(db *sqlx.DB) domain.EntryRepository {
	return sqliteEntryRepository{db: db}
}

func (s sqliteEntryRepository) GetById(id int64, withJson bool) (domain.Entry, error) {
	entry := domain.Entry{}
	selectedRows := "*"
	if !withJson {
		selectedRows = "id, provider_id, url, title, published_at, author, fetched_at"
	}
	err := s.db.Get(&entry, "SELECT "+selectedRows+" FROM entries WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return entry, err
	}
	return entry, nil
}

func (s sqliteEntryRepository) GetAll() (*[]domain.Entry, error) {
	panic("implement me")
}

func (s sqliteEntryRepository) Insert(entry domain.Entry) error {
	res, err := s.db.NamedExec("INSERT INTO entries (provider_id, url, title, published_at, author, fetched_at, json) VALUES (:provider_id, :url, :title, :published_at, :author, :fetched_at, :json)", entry)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	} // TODO: Throw error if no rows affected.
	return err
}

func (s sqliteEntryRepository) BulkInsert(entries []domain.Entry) error {
	res, err := s.db.NamedExec("INSERT INTO entries (provider_id, url, title, published_at, author, fetched_at, json) VALUES (:provider_id, :url, :title, :published_at, :author, :fetched_at, :json)", entries)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	} // TODO: Throw error if no rows affected.
	return err
}

func (s sqliteEntryRepository) Update(entry domain.Entry) error {
	res, err := s.db.NamedExec("UPDATE entries SET url = :url, title = :title, published_at = :published_at, author = :author, fetched_at = :fetched_at, json = :json WHERE id = :id", entry)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	} // TODO: Throw error if no rows affected.
	return err
}

func (s sqliteEntryRepository) Delete(id int64) error {
	panic("implement me")
}

func (e sqliteEntryRepository) GetByQuery(query string, dateFrom, dateUntil, providerId, limit, offset int64, withJson, onlyIdAndUrl bool) ([]domain.Entry, error) {
	var entries []domain.Entry
	var sqlQuery string

	whereClauses := []string{}
	if providerId != -1 {
		whereClauses = append(whereClauses, "provider_id = :providerId")
	}
	if dateFrom != -1 && dateUntil != -1 {
		whereClauses = append(whereClauses, "(published_at BETWEEN :dateFrom AND :dateUntil)")
	}
	if query != "" {
		whereClauses = append(whereClauses, "(title LIKE :query)")
	}
	if onlyIdAndUrl {
		sqlQuery = "SELECT id, url FROM entries"
	} else if withJson {
		sqlQuery = "SELECT * FROM entries"
	} else {
		sqlQuery = "SELECT id, provider_id, url, title, published_at, author, fetched_at FROM entries"
	}
	if len(whereClauses) != 0 {
		sqlQuery += " WHERE "
	}
	sqlQuery += strings.Join(whereClauses, " AND ")
	sqlQuery += " ORDER BY published_at DESC LIMIT :limit OFFSET :offset"

	rows, err := e.db.NamedQuery(sqlQuery, map[string]interface{}{"providerId": providerId, "query": ("%" + query + "%"), "dateFrom": dateFrom, "dateUntil": dateUntil, "limit": limit, "offset": offset})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		entry := domain.Entry{}
		err := rows.StructScan(&entry)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
