package queryhand

import (
	"context"
	"encoding/json"

	"github.com/coocood/freecache"
)

var cache *freecache.Cache

func SetCache(c *freecache.Cache) { cache = c }

func QueryLocalCache(ctx context.Context, param any, res any) error {
	key := ""
	value, err := cache.Get([]byte(key))
	if err != nil {
		return err
	}
	if len(value) == 0 {

	}
	err = json.Unmarshal(value, res)
	if err != nil {
		return err
	}
	cache.Set([]byte(key), value, 5)
	return nil
}

func QueryRedisCache(ctx context.Context, param any, res any) error {
	return nil
}

func QueryMySQLDB(ctx context.Context, param any, res any) error {
	return nil
}
