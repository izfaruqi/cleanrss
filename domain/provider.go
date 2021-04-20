package domain

type Provider struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Url       string `json:"url" db:"url"`
	ParserId  int64  `json:"parserId" db:"parserId"`
	IsDeleted bool   `json:"is_deleted" db:"isDeleted"`
}

type ProviderUsecase interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider Provider) error
	Update(id int64, provider Provider) error
	Delete(id int64) error
}

type ProviderRepository interface {
	GetById(id int64) (Provider, error)
	GetAll() (*[]Provider, error)
	Insert(provider Provider) error
	Update(id int64, provider Provider) error
	Delete(id int64) error
}
