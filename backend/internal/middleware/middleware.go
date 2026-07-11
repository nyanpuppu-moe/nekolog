package middleware

import (
	"net/http"

	"nekolog/internal/utils"
	"nekolog/internal/web"
)

func AuthRequired(c *web.Context) {
	userID := c.SessionGet("user_id")

	if userID == nil {
		c.JSON(
			http.StatusUnauthorized,
			utils.Object{
				"error": "Faild auth, please account login",
			},
		)
		c.Abort()
		return
	}

	c.SessionSet("user_id", userID)

	c.Next()
}
