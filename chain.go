package queryhand

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found record")

type Query interface {
	Next() bool
	Scan(ctx context.Context, param any, res any) error
	Exec(ctx context.Context, param any, res any) error
}

type Handler func(ctx context.Context, param any, res any) error

type ExampleQuery struct {
	Len      uint32    //处理链路长度
	Current  uint32    //当前处理器位移
	Handlers []Handler //处理器列表
}

func NewQuery(handlers ...Handler) *ExampleQuery {
	return &ExampleQuery{
		Len:      uint32(len(handlers)),
		Current:  0,
		Handlers: handlers,
	}
}

func (h *ExampleQuery) Next() bool {
	// 判断是否还存在处理器未执行
	return h.Current == h.Len
}

func (h *ExampleQuery) Scan(ctx context.Context, param, res any) error {
	defer func() { h.Current++ }()
	return h.Handlers[h.Current](ctx, param, res)
}

func (query *ExampleQuery) Exec(ctx context.Context, param any, res any) error {
	for query.Next() {
		var param, res any
		err := query.Scan(context.Background(), &param, &res)
		if err != nil {
			return err
		}
		if res != nil {
			// 已获取到结果，直接返回
			return nil
		}
	}
	return ErrNotFound
}
