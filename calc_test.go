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

var (
	POSITION_DONGDAEMOON       Coordinate = Coordinate{"신평화패션타운", 37.56946547219505, 127.01168167917474}
	POSITION_SEOULSTATION      Coordinate = Coordinate{"서울역", 37.55463197345425, 126.97055261318205}
	POSITION_HANYANGUNIV       Coordinate = Coordinate{"한양대학교", 37.55772638076536, 127.04531599177925}
	POSITION_GANGNAMHYUNDAIAPT Coordinate = Coordinate{"강남현대아파트", 37.52635073205962, 127.02442463468053}
	POSITION_GALLERIAFOREST    Coordinate = Coordinate{"서울숲갤러리아포레", 37.54593339811107, 127.04240066671743}
)

var mockRequest CalculateRequest = CalculateRequest{
	Drivers: []Driver{
		Driver{"김경식", Coordinate{"d1", 0, 0}, 100},
		Driver{"김진희", Coordinate{"d2", 100, 100}, 100},
	},
	Stuffs: []Stuff{
		Stuff{"후드티", "김기범", POSITION_DONGDAEMOON, "김정현", POSITION_GANGNAMHYUNDAIAPT},
		Stuff{"청바지", "김기범", POSITION_DONGDAEMOON, "조현재", POSITION_HANYANGUNIV},
		Stuff{"원단", "김기범", POSITION_DONGDAEMOON, "임동영", POSITION_GALLERIAFOREST},
	},
}

var (
	kmeanResult       clusters.Clusters
	pairClusterResult []PairCluster
	graphResult       []DistanceGraph
	driverResult      map[string]*DistanceGraph
)

func TestKmean(t *testing.T) {
	fmt.Println(aurora.Bold(aurora.BgMagenta("Input")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, mockRequest)), nil)))

	coords := mockRequest.Stuffs.Coordinates("center")
	cs, err := GetKmeanCluster(coords, 2)
	assert.NoError(t, err)

	fmt.Println(aurora.Bold(aurora.BgMagenta("Center Result")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, coords)), nil)))
	fmt.Println(aurora.Bold(aurora.BgMagenta("Kmean Result")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, cs)), nil)))

	kmeanResult = cs
}

func TestConvertToPair(t *testing.T) {
	pairClusterResult = convertCenterToPairs(mockRequest, kmeanResult)
	fmt.Println(aurora.Bold(aurora.BgMagenta("Result to Pair")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, pairClusterResult)), nil)))
}

func TestExtractPairsToCoord(t *testing.T) {
	for _, p := range pairClusterResult {
		res := p.Pairs.Coordinates("all")
		fmt.Println(aurora.Bold(aurora.BgMagenta("Pair to Coordinates")))
		fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, res)), nil)))
	}
}

func prepareApiTest() error {
	InitMapClient()
	if c, err := ReadConfig("config.json", nil); err != nil {
		return err
	} else {
		Config = c
		return nil
	}
}

func TestMapApi(t *testing.T) {
	return
	assert.NoError(t, prepareApiTest())
	d, err := getRoadDistance(Coordinate{"start", 127.33, 37.5}, Coordinate{"goal", 127.55, 36}, nil)
	assert.NoError(t, err)
	if err != nil {
		fmt.Print(aurora.Bold(aurora.BgMagenta("Distance")), *d, "m\n")
	}
}

func TestDistanceGraph(t *testing.T) {
	assert.NoError(t, prepareApiTest())
	graphResult = MakeDistanceGraph(pairClusterResult, func(format string, args ...interface{}) {
		//assert.Failf(t, "MakeDistanceGraph() error : ", format, args)
		fmt.Printf(format, args)
	}, nil)

	for k, g := range graphResult {
		fmt.Println(aurora.Bold(aurora.BgMagenta("Distance graph")), fmt.Sprintf("#%v", k))
		fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, g)), nil)))
	}
}

func TestAssignDriver(t *testing.T) {
	driverResult = AssignDriverToGraphs(graphResult, mockRequest.Drivers)
	fmt.Println(aurora.Bold(aurora.BgMagenta("Assign driver")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, driverResult)), nil)))
}

func TestFindActions(t *testing.T) {
	fmt.Println(aurora.Bold(aurora.BgMagenta("Action")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, FindActions(driverResult, mockRequest))), nil)))

}

func TestCalculateActions(t *testing.T) {
	r, err := calculateActions(mockRequest)
	assert.NoError(t, err)
	fmt.Println(aurora.Bold(aurora.BgMagenta("Final")))
	fmt.Println(string(pretty.Color(pretty.Pretty(mustMarshal(t, r)), nil)))
}

func mustMarshal(t *testing.T, i interface{}) []byte {
	j, err := json.Marshal(i)
	assert.NoError(t, err)
	return j
}
