package main

import (
	"github.com/labstack/echo/v4"

	"net/http"
)

// @summary API 자격증명 토큰을 발급합니다
// @description 액세스 키를 제출받고, 유효한 액세스키라면 API 자격증명 토큰을 발급합니다
// @id issue-token
// @accept json
// @produce json
// @param key body string true "액세스 키"
// @success 200 {body} string "API 자격증명 토큰"
// @failure 400
// @failure 401
// @router /token [POST]
func rIssueToken(c echo.Context) error {
	req := map[string]interface{}{}
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

// @summary 배송 경로를 계산합니다.
// @description 제출된 배송자, 배송 물품 정보로 배송 경로를 계산합니다.
// @id calculate-path
// @accept json
// @produce json
// @securityDefinitions.apikey API-TOKEN
// @param API-TOKEN header string true "API 자격증명 코드"
// @param info body CalculateRequest true "배송자와 배송 물품 정보"
// @success 200 {object} CalculateResult "배송 경로"
// @failure 400
// @failure 401
// @router /path [POST]
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
