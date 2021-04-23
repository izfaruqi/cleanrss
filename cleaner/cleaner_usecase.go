package cleaner

import "cleanrss/domain"

type cleanerUsecase struct {
	r domain.CleanerRepository
}

func NewCleanerUsecase(r domain.CleanerRepository) domain.CleanerUsecase {
	return cleanerUsecase{r: r}
}

func (c cleanerUsecase) GetById(id int64) (domain.Cleaner, error) {
	return c.r.GetById(id)
}

func (c cleanerUsecase) GetAll() (*[]domain.Cleaner, error) {
	return c.r.GetAll()
}

func (c cleanerUsecase) Insert(cleaner *domain.Cleaner) error {
	return c.r.Insert(cleaner)
}

func (c cleanerUsecase) Update(cleaner domain.Cleaner) error {
	return c.r.Update(cleaner)
}

func (c cleanerUsecase) Delete(id int64) error {
	return c.r.Delete(id)
}
