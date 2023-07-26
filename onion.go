package queryhand

import "math"

const abortIndex int8 = math.MaxInt8 >> 1

type Handler func(ctx *Context) error

type Context struct {
	index    int8
	handlers []Handler
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}
