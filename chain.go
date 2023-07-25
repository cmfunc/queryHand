package queryhand

import "github.com/pkg/errors"

var ErrNotFound = errors.New("not found record")

type Handler interface {
	GetFromLocal(param, resp interface{}) error
	GetFromCache(param, resp interface{}) error
	GetFromDB(param, resp interface{}) error
}

// Query is applay Handler interface to exec
// query from localcache, distributed cache, db.
func Query(h Handler, param, resp interface{}) error {
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
