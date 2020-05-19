package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	//KeyRequest is form of issueing token request.
	KeyRequest struct {
		Key string `json:"key" validate:"required"`
	}

	//TokenResponse is form of issueing token response.
	TokenResponse struct {
		Token string `json:"token"`
	}

	//CalculateRequest is structure for api request.
	CalculateRequest struct {
		Drivers []Driver `json:"drivers" validate:"required"` //Current available drivers data
		Stuffs  []Stuff  `json:"stuffs" validate:"required"`  //Current available stuffs data
	}

	//CalculateResult is structure for api response.
	CalculateResult struct {
		Actions map[string][]DriverAction `json:"actions"` //key is Driver's id
	}
)

// @summary API 자격증명 토큰을 발급합니다.
// @description 액세스 키를 제출받고, 유효한 액세스키라면 API 자격증명 토큰을 발급합니다
// @id issue-token
// @accept json
// @produce json
// @param key body KeyRequest true "액세스 키"
// @success 200 {object} TokenResponse "유효한 액세스 키를 제출했을때 API 자격증명 토큰이 반환됩니다."
// @failure 400 "요청이 정해진 형식에 부합하지 않습니다."
// @failure 401 "액세스 키가 유효하지 않습니다."
// @router /token [POST]
func rIssueToken(c echo.Context) error {
	req := KeyRequest{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(req); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if searchSlice(req.Key, &Config.AccessKey) {
		return c.JSON(http.StatusOK, TokenResponse{newToken().Token})
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
// @success 200 {object} CalculateResult "요청에 대한 배송 경로가 반환됩니다."
// @failure 400 "요청이 정해진 형식에 부합히지 않습니다."
// @failure 401 "액세스 키가 유효하지 않습니다."
// @failure 500 "처리 중 서버 내부 오류가 발생했습니다."
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

	if d, err := mockCalculateActions(req); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else {
		return c.JSON(http.StatusOK, d)
	}

}
