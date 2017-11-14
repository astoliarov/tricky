package dao

import (
	"github.com/cornelk/hashmap"
	"tricky/lib/domain"
	"unsafe"
)

type InMemoryDAO struct {
	hMap hashmap.HashMap
}

func NewInMemoryRepository() *InMemoryDAO {
	repo := InMemoryDAO{}
	return &repo
}

func (r *InMemoryDAO) Add(rule *domain.RedirectionRule) error {
	r.hMap.Set(rule.Key, unsafe.Pointer(&rule))
	return nil
}

func (r *InMemoryDAO) Delete(key string) (*domain.RedirectionRule, error) {
	rule, _ := r.GetByKey(key)
	if rule != nil {
		r.hMap.Del(key)
	}
	return rule, nil
}

func (r *InMemoryDAO) GetByKey(key string) (*domain.RedirectionRule, error) {
	res, ok := r.hMap.Get(key)
	if !ok {
		return nil, nil
	}

	return *(**domain.RedirectionRule)(res), nil
}

func (r *InMemoryDAO) GetAll() ([]*domain.RedirectionRule, error) {
	var rules []*domain.RedirectionRule

	for value := range r.hMap.Iter() {
		rules = append(rules, *(**domain.RedirectionRule)(value.Value))
	}

	return rules, nil
}
