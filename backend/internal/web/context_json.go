package web

import "encoding/json"

func (c *Context) BindJSON(obj any) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

func (c *Context) JSON(code int, obj any) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(obj)
}
