package provider

import (
	"cleanrss/domain"
)

type providerUsecase struct {
	r domain.ProviderRepository
}

func NewProviderUsecase(r domain.ProviderRepository) domain.ProviderUsecase {
	return &providerUsecase{
		r: r,
	}
}

func (m providerUsecase) GetById(id int64) (domain.Provider, error) {
	return m.r.GetById(id)
}

func (m providerUsecase) GetAll() (*[]domain.Provider, error) {
	return m.r.GetAll()
}

func (m providerUsecase) Insert(provider *domain.Provider) error {
	return m.r.Insert(provider)
}

func (m providerUsecase) Update(provider domain.Provider) error {
	return m.r.Update(provider)
}

func (m providerUsecase) Delete(id int64) error {
	return m.r.Delete(id)
}
