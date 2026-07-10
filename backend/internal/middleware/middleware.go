package middleware

import (
	"nekolog/internal/web"
	"net/http"
)

func AuthRequired(c *web.Context) {
	userID := c.SessionGet("user_id")

	if userID == nil {
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "Faild auth, please account login",
			},
		)
		c.Abort()
		return
	}

	c.SessionSet("user_id", userID)

	c.Next()
}
