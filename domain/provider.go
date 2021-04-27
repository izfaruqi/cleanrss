package domain

type Provider struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Url       string `json:"url" db:"url"`
	ParserId  int64  `json:"parserId" db:"parser_id"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
}

type ProviderUsecase interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider *Provider) error
	Update(provider Provider) error
	Delete(id int64) error
}

type ProviderRepository interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider *Provider) error
	Update(provider Provider) error
	Delete(id int64) error
}
