package main

import (
	"github.com/labstack/echo/v4"

	"net/http"
)

func rIssueToken(c echo.Context) error {
	req := map[interface{}]interface{}{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	var reqkey string
	if val, res := req["key"].(string); !res {
		c.Logger().Error("Key retrieve fail")
		return c.NoContent(http.StatusBadRequest)
	} else {
		reqkey = val
	}

	if searchSlice(reqkey, &validAccessKey) {
		return c.JSON(http.StatusOK, map[string]string{"token": newToken().Token})
	} else {
		return c.NoContent(http.StatusUnauthorized)
	}
}

func rCalculatePath(c echo.Context) error {
	if tok := c.Request().Header.Get("API-TOKEN"); tok == "" {
		return c.NoContent(http.StatusBadRequest)
	} else {
		if !validateToken(tok) {
			return c.NoContent(http.StatusUnauthorized)
		}
	}

	req := CalculateRequest{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, calculateActions(req))
}
