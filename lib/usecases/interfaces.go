package usecases

import "tricky/lib/domain"

//go:generate mockgen -destination=rules_dao_mock.go -package=usecases tricky/lib/usecases IRulesDAO

type IRulesDAO interface {
	Delete(key string) (*domain.RedirectionRule, error)
	Add(rule *domain.RedirectionRule) error
	GetByKey(key string) (*domain.RedirectionRule, error)
	GetAll() ([]*domain.RedirectionRule, error)
}

type ILogger interface {
	Info(message string)
	Debug(message string)
}
