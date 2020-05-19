package main

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/muesli/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/pretty"
	"testing"
)

var mockCoordinate []*Coordinate = []*Coordinate{
	&Coordinate{"test1", 0.0, 0.0},
	&Coordinate{"test2", 1.0, 1.0},
	&Coordinate{"test3", -1.0, 1.0},
	&Coordinate{"test4", 1.0, -1.0},
	&Coordinate{"test5", -1.0, -1.0},
	&Coordinate{"test6", 11.0, 11.0},
	&Coordinate{"test7", 9.0, 11.0},
	&Coordinate{"test8", 11.0, 9.0},
	&Coordinate{"test9", 9.0, 9.0},
}

var mockRequest CalculateRequest = CalculateRequest{
	Drivers: []Driver{
		Driver{"John", Coordinate{"d1", 0, 0}, 100},
		Driver{"John", Coordinate{"d2", 100, 100}, 100},
	},
	Stuffs: []Stuff{
		Stuff{"apple", "Kim", Coordinate{"s1s", 12.73, 14.94}, "Ha", Coordinate{"s1r", 72.12, 9.43}},
		Stuff{"grape", "Jang", Coordinate{"s2s", 78.16, 58.54}, "Won", Coordinate{"s2r", 78.20, 45.62}},
		Stuff{"chicken", "Park", Coordinate{"s3s", 97.46, 41.98}, "Go", Coordinate{"s3r", 36.43, 93.96}},
		Stuff{"pizza", "Choi", Coordinate{"s4s", 43.48, 85.71}, "Hwang", Coordinate{"s4r", 85.21, 34.98}},
		Stuff{"soup", "Min", Coordinate{"s5s", 79.35, 91.31}, "Nho", Coordinate{"s5r", 20.38, 75.61}},
		Stuff{"hamberger", "Hong", Coordinate{"s6s", 24.89, 35.96}, "Jeon", Coordinate{"s6r", 65.73, 23.05}},
		Stuff{"rice", "Mo", Coordinate{"s7s", 59.21, 25.48}, "Moon", Coordinate{"s7r", 25.23, 0.65}},
		Stuff{"ramen", "Woo", Coordinate{"s8s", 23.65, 11.74}, "Son", Coordinate{"s8r", 98.64, 21.28}},
	},
}

var kmeanResult clusters.Clusters

func TestKmean(t *testing.T) {
	coords := extractCoordinate("center", &mockRequest.Stuffs)
	cs, err := GetKmeanCluster(coords, 2)
	assert.NoError(t, err)

	fmt.Println(aurora.Bold(aurora.BgMagenta("Input")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, coords)), nil)))
	fmt.Println(aurora.Bold(aurora.BgMagenta("Result")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, cs)), nil)))

	kmeanResult = cs
}

func estConvertCluster(t *testing.T) {
	fmt.Println(aurora.Bold(aurora.BgMagenta("Convert")))
	//fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, convertClusters(mockRequest, kmeanResult))), nil)))
}

func mustMarshal(t *testing.T, i interface{}) []byte {
	j, err := json.Marshal(i)
	assert.NoError(t, err)
	return j
}
