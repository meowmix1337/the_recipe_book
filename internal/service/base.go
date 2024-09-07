package service

import (
	"strings"

	"github.com/meowmix1337/go-core/cache"
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/segmentio/ksuid"
)

type BaseService struct {
	Config config.Config
	Cache  cache.Cache
}

func NewBaseService(cfg config.Config, cache cache.Cache) *BaseService {
	return &BaseService{
		Config: cfg,
		Cache:  cache,
	}
}

func (s *BaseService) GenerateUUIDHash(prefix string) string {
	id := ksuid.New()

	prefix = strings.TrimSpace(prefix)
	prefix = strings.ReplaceAll(prefix, "_", "")

	return prefix + "_" + id.String()
}
