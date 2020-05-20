package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

var globalClient *resty.Client

type NaverResponse struct {
	Code            int                    `json:"code"`
	Message         string                 `json:"message"`
	CurrentDateTime string                 `json:"currentDateTime"`
	Route           map[string]interface{} `json:"route"`
}

func MakeDistanceGraph() {

}

func InitClient() {
	globalClient = resty.New()
}

//https://apidocs.ncloud.com/ko/ai-naver/maps_directions/driving/
func getRoadDistance(start, goal Coordinate) (d *float64, err error) {
	resp, err := globalClient.R().
		SetHeader("X-NCP-APIGW-API-KEY-ID", Config.NaverClientId).
		SetHeader("X-NCP-APIGW-API-KEY", Config.NaverClientSecret).
		SetQueryParams(map[string]string{
			"start": coordToNaverFormat(start),
			"goal":  coordToNaverFormat(goal),
		}).
		Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Api provider return http response code %d (200 expected)", resp.StatusCode())
	}

	res := NaverResponse{}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("Api provider return result code %d (0 expected)", res.Code)
	}

	defer func() {
		if r := recover(); r != nil {
			d = nil
			err = fmt.Errorf("Api response parsing error : %s", r)
		}
	}()
	for _, v := range res.Route {
		v = v.([]interface{})[0]
		v = v.(map[string]interface{})["summary"]
		v = v.(map[string]interface{})["distance"]
		var vint float64 = v.(float64)
		d = &vint
		err = nil
	}
	return
}

func coordToNaverFormat(c Coordinate) string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Long)
}
