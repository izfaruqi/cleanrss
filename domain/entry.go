package domain

type Entry struct {
	Id          int64  `json:"id" db:"id"`
	ProviderId  int64  `json:"providerId" db:"providerId"`
	Url         string `json:"url" db:"url"`
	Title       string `json:"title" db:"title"`
	PublishedAt int64  `json:"publishedAt" db:"publishedAt"`
	Author      string `json:"author" db:"author"`
	FetchedAt   int64  `json:"fetchedAt" db:"fetchedAt"`
	Json        string `json:"json" db:"json"`
}

type EntryUsecase interface {
	GetById(id int64, withJson bool) (Entry, error)
	GetAll(withJson bool) (*[]Entry, error)
	GetByQuery(query string, dateFrom int64, dateUntil int64, providerId int64, limit int64, offset int64, withJson bool) (*[]Entry, error)
	TriggerRefresh(providerIds []int64) error
}

type EntryRepository interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider Provider) error
	Update(id int64, provider Provider) error
	Delete(id int64) error
}

type EntrySource interface {
	GetByProviderId(providerId int64) (*[]Entry, error)
}
