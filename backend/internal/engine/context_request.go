package engine

import "mime/multipart"

func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	err := c.Request.ParseMultipartForm(32 << 20)

	if err != nil {
		return nil, err
	}

	file, header, err := c.Request.FormFile(key)

	if err != nil {
		return nil, err
	}

	file.Close()

	return header, nil
}
