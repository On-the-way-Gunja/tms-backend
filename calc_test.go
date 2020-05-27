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
	POSITION_1  Coordinate = Coordinate{"중앙대학교병원", 37.50483772232202, 126.95498742774899}
	POSITION_2  Coordinate = Coordinate{"노량진역", 37.51401333884956, 126.9420872919965}
	POSITION_3  Coordinate = Coordinate{"서울역", 37.554586923571534, 126.9705526309084}
	POSITION_4  Coordinate = Coordinate{"고려대학교서울캠퍼스", 37.58994135115182, 127.03182770693546}
	POSITION_5  Coordinate = Coordinate{"신평화패션타운", 37.56948350546714, 127.01154584848007}
	POSITION_6  Coordinate = Coordinate{"경희대학교서울캠퍼스", 37.5967813287537, 127.0528925020005}
	POSITION_7  Coordinate = Coordinate{"서울과학기술대학교", 37.63212235619019, 127.0777506016829}
	POSITION_8  Coordinate = Coordinate{"목동중학교", 37.52073335886291, 126.8716643597581}
	POSITION_9  Coordinate = Coordinate{"대성디폴리스지식산업센터", 37.48008564591763, 126.87697987125098}
	POSITION_10 Coordinate = Coordinate{"성산e편한세상2차아파트", 37.570067820692714, 126.90535727247682}
	POSITION_11 Coordinate = Coordinate{"덕수궁롯데캐슬아파트", 37.5639752578396, 126.97033388202121}
	POSITION_12 Coordinate = Coordinate{"건영캐스빌아파트", 37.60290247574357, 127.09149079480794}
	POSITION_13 Coordinate = Coordinate{"하남코스트코", 37.54987128494775, 127.1945417296819}
	POSITION_14 Coordinate = Coordinate{"한국외국어대학교서울캠퍼스", 37.59773359335199, 127.05878151886373}
)

var mockRequest CalculateRequest = CalculateRequest{
	Drivers: []Driver{
		Driver{"김승우", POSITION_1, 100},
		Driver{"김인섭", POSITION_2, 100},
		Driver{"김정현", POSITION_3, 100},
	},
	Stuffs: []Stuff{
		Stuff{"청바지", "안권우", POSITION_4, "정연길", POSITION_5},
		Stuff{"후드티", "조우현", POSITION_6, "정연길", POSITION_5},
		Stuff{"티셔츠", "문준수", POSITION_7, "정연길", POSITION_5},
		Stuff{"유압펌프", "황동엽", POSITION_8, "박준규", POSITION_9},
		Stuff{"LED매트릭스", "이동찬", POSITION_10, "박준규", POSITION_9},
		Stuff{"모니터", "조현재", POSITION_11, "박준규", POSITION_9},
		Stuff{"책상", "박세현", POSITION_12, "임동영", POSITION_13},
		Stuff{"장롱", "전민수", POSITION_14, "임동영", POSITION_13},
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
	cs, err := GetKmeanCluster(coords, len(mockRequest.Drivers))
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
	InitMap()
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
