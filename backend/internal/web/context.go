package web

import (
	"net/http"
)

type Context struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	store     *SessionStore
	session   Session
	sessionID string

	handlers []HandlerFunc
	index    int
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}
