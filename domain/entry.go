package domain

type Entry struct {
	Id          int64  `json:"id" db:"id"`
	ProviderId  int64  `json:"providerId" db:"provider_id"`
	Url         string `json:"url" db:"url"`
	Title       string `json:"title" db:"title"`
	PublishedAt int64  `json:"publishedAt" db:"published_at"`
	Author      string `json:"author" db:"author"`
	FetchedAt   int64  `json:"fetchedAt" db:"fetched_at"`
	Json        string `json:"json" db:"json"`
}

type EntryUsecase interface {
	GetById(id int64, withJson bool) (Entry, error)
	GetAll(withJson bool) (*[]Entry, error)
	GetByQuery(query string, dateFrom int64, dateUntil int64, providerId int64, limit int64, offset int64, withJson bool) (*[]Entry, error)
	TriggerRefresh(providerId int64) error
	TriggerRefreshAll() error
}

type EntryRepository interface {
	GetById(id int64, withJson bool) (Entry, error)
	GetAll() (*[]Entry, error)
	Insert(provider Entry) error
	Update(provider Entry) error
	Delete(id int64) error
	GetByQuery(query string, dateFrom int64, dateUntil int64, providerId int64, limit int64, offset int64, withJson bool) (*[]Entry, error)
	GetUrlById(id int64) (string, error)
	BulkInsert(entries []Entry) error
}

type WebExtEntryRepository interface {
	GetRawEntriesByProviderId(providerId int64) ([]Entry, int, error)
}
