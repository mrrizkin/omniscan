package provider

import (
	"github.com/dgraph-io/ristretto/v2"

	"github.com/mrrizkin/omniscan/config"
)

type Ristretto struct {
	cache  *ristretto.Cache[string, interface{}]
	config *config.App
}

func NewRessetto(config *config.App) (*Ristretto, error) {
	ristrettoConfig := ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	}
	cache, err := ristretto.NewCache(&ristrettoConfig)
	if err != nil {
		return nil, err
	}

	return &Ristretto{
		cache:  cache,
		config: config,
	}, nil
}

func (c *Ristretto) Set(key string, value interface{}) {
	c.cache.SetWithTTL(key, value, 0, c.config.CacheTTLSecond())
}

func (c *Ristretto) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Ristretto) Delete(key string) {
	c.cache.Del(key)
}

func (c *Ristretto) Close() error {
	c.cache.Close()

	return nil
}
