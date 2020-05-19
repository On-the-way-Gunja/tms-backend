package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func MakeDistanceGraph() {

}

//https://apidocs.ncloud.com/ko/ai-naver/maps_directions/driving/
func getRoadDistance() {
	client := resty.New()
	resp, err := client.R().
		SetHeader("X-NCP-APIGW-API-KEY-ID", Config.NaverClientId).
		SetHeader("X-NCP-APIGW-API-KEY", Config.NaverClientSecret).
		SetQueryParams(map[string]string{
			"start": "",
			"goal":  "",
		}).
		Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
}

func coordToNaverFormat(c Coordinate) string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Long)
}
