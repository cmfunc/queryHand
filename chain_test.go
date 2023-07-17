package queryhand

import (
	"context"
	"testing"
)

func QueryLocalCache(ctx context.Context, param any, res any) error {
	return nil
}

func QueryRedisCache(ctx context.Context, param any, res any) error {
	return nil
}

func QueryMySQLDB(ctx context.Context, param any, res any) error {
	return nil
}

func TestQuery(t *testing.T) {
	query := NewQuery(QueryLocalCache, QueryRedisCache, QueryMySQLDB)
	for query.Next() {
		var param, res any
		err := query.Scan(context.Background(), &param, &res)
		if err != nil {
			t.Error(err)
			return
		}
		if res != nil {
			// 已获取到结果，直接返回
			return
		}
	}
}
