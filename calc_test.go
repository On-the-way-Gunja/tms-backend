package main

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/pretty"
	"testing"
)

var CoordinateTestset []*Coordinate = []*Coordinate{
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

func TestKmean(t *testing.T) {
	cs, err := GetKmeanCluster(CoordinateTestset, 2)
	assert.NoError(t, err)

	fmt.Println(aurora.Bold(aurora.BgMagenta("Input")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, CoordinateTestset)), nil)))
	fmt.Println(aurora.Bold(aurora.BgMagenta("Result")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, cs)), nil)))
}

func mustMarshal(t *testing.T, i interface{}) []byte {
	j, err := json.Marshal(i)
	assert.NoError(t, err)
	return j
}
