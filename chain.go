package queryhand

import (
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("not found record")

type CacheType int

const (
	_ CacheType = iota
	CacheTypeString
	CacheTypeList
	CacheTypeHash
	CacheTypeSet
	CacheTypeZset
)

type Param interface {
	LocalKey() string
	CacheKeyAndType() (string, CacheType) //返回Redis缓存的key以及存储数据类型
	QuerySQL() string                     //返回查询SQL
}

type QueryHandler interface {
	GetFromLocal(param Param, resp interface{}) error
	GetFromCache(param Param, resp interface{}) error
	GetFromDB(param Param, resp interface{}) error
}

type ExecQueryLocal func(key string, resp interface{}) error
type ExecQueryCache func(key string, valType CacheType, resp interface{}) error
type ExecQueryDB func(query string, resp interface{}) error

func NewQueryhandler(execLocal ExecQueryLocal) *queryHandler {
	return &queryHandler{}
}

type queryHandler struct {
	execLocalFunc ExecQueryLocal
	execCacheFunc ExecQueryCache
	execDBFunc    ExecQueryDB
}

func (h *queryHandler) GetFromLocal(param Param, resp interface{}) error {
	key := param.LocalKey()
	err := h.execLocalFunc(key, resp)
	if err != nil {
		return errors.Wrapf(err, "queryHandler GetFromLocal param: %+v resp:%+v", param, resp)
	}
	return nil
}
func (h *queryHandler) GetFromCache(param Param, resp interface{}) error {
	key, valType := param.CacheKeyAndType()
	err := h.execCacheFunc(key, valType, resp)
	if err != nil {
		return errors.Wrapf(err, "queryHandler GetFromCache param: %+v resp:%+v", param, resp)
	}
	return nil
}
func (h *queryHandler) GetFromDB(param Param, resp interface{}) error {
	query := param.QuerySQL()
	err := h.execDBFunc(query, resp)
	if err != nil {
		return errors.Wrapf(err, "queryHandler GetFromCache param: %+v resp:%+v", param, resp)
	}
	return nil
}

// Query is applay Handler interface to exec
// query from localcache, distributed cache, db.
func Query(h QueryHandler, param Param, resp interface{}) error {
	// query from local cache
	err := h.GetFromLocal(param, resp)
	if err != nil {
		return errors.Wrapf(err, "queryhand Query h:%+v", h)
	}
	if resp != nil {
		return nil
	}
	// query from distruted cache
	err = h.GetFromCache(param, resp)
	if err != nil {
		return errors.Wrapf(err, "queryhand Query h:%+v", h)
	}
	if resp != nil {
		return nil
	}
	// query from db
	err = h.GetFromDB(param, resp)
	if err != nil {
		return errors.Wrapf(err, "queryhand Query h:%+v", h)
	}
	if resp != nil {
		return nil
	}
	return ErrNotFound
}
