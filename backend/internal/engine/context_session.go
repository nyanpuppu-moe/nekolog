package engine

import "net/http"

func (c *Context) SessionGet(key string) any {
	return c.session[key]
}

func (c *Context) SessionSet(key string, value any) {
	if c.session == nil {
		c.session = make(map[string]any)
	}
	c.session[key] = value
}

func (c *Context) SessionSave() error {
	if c.session == nil {
		return nil
	}

	if c.sessionID == "" {
		id, err := c.store.GenerateID()
		if err != nil {
			return err
		}
		c.sessionID = id

		cookie := &http.Cookie{
			Name:     SessionCookieName,
			Value:    c.sessionID,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
	}

	c.store.Set(c.sessionID, c.session)

	return nil
}

func (c *Context) SessionClear() {
	c.session = make(Session)

	if c.sessionID != "" {
		c.store.Delete(c.sessionID)
	}

	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, cookie)

	c.sessionID = ""
}
