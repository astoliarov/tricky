package usecases

import (
	"net/url"

	"tricky/lib/domain"
)

type RulesUseCase struct {
	rulesDAO IRulesDAO
	//logger   ILogger
}

func (ri *RulesUseCase) FindUrl(keys []string) (*url.URL, error) {
	var rule *domain.RedirectionRule
	for i := range keys {
		key := keys[i]

		foundedRule, err := ri.rulesDAO.GetByKey(key)
		if err != nil {
			// TODO: create abstract storage error
			return nil, err
		}
		if foundedRule != nil {
			rule = foundedRule
			break
		}
	}
	if rule == nil {
		return nil, nil
	}

	return rule.Url, nil
}

func (ri *RulesUseCase) AddRule(key string, url *url.URL) error {
	rule := domain.RedirectionRule{key, url}
	return ri.rulesDAO.Add(&rule)
}

func (ri *RulesUseCase) DeleteRule(key string) (*domain.RedirectionRule, error) {
	return ri.rulesDAO.Delete(key)
}

func (ri *RulesUseCase) GetAllRules() ([]*domain.RedirectionRule, error) {
	return ri.rulesDAO.GetAll()
}

func NewRulesUseCase(dao IRulesDAO) *RulesUseCase {
	useCase := RulesUseCase{rulesDAO: dao}
	return &useCase
}
