package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	//"github.com/muesli/clusters"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/yourbasic/graph"
	"math"
)

var globalClient *resty.Client
var cacheDb *badger.DB

func InitMap() (func(), error) {
	globalClient = resty.New()
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, err
	} else {
		cacheDb = db
		return func() {
			db.Close()
		}, nil
	}
}

func callDistanceApi(start, goal Coordinate) (*[]byte, *NaverResponse, error) {
	currentId := start.Id + "-" + goal.Id
	var body []byte
	iscached := false

	err := cacheDb.View(func(txn *badger.Txn) error {
		item, errkey := txn.Get([]byte(currentId))
		if errkey == nil {
			if Logger != nil {
				Logger.Tracef("[callDistanceApi] Cache hit!: %s", currentId)
			}
			iscached = true
			err := item.Value(func(val []byte) error {
				body = append([]byte{}, val...)
				return nil
			})
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !iscached {
		resp, err := globalClient.R().
			SetHeader("X-NCP-APIGW-API-KEY-ID", Config.NaverClientId).
			SetHeader("X-NCP-APIGW-API-KEY", Config.NaverClientSecret).
			SetQueryParams(map[string]string{
				"start": coordToNaverFormat(start),
				"goal":  coordToNaverFormat(goal),
			}).
			Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
		if err != nil {
			return nil, nil, err
		}
		body = resp.Body()
		if resp.StatusCode() != 200 {
			return &body, nil, fmt.Errorf("Api provider return http response code %d (200 expected)", resp.StatusCode())
		} else {
			err := cacheDb.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte(currentId), body)
			})
			if err != nil {
				return &body, nil, err
			}
		}
	}
	res := NaverResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return &body, nil, err
	}

	return &body, &res, nil
}

func getRoadDistance(start, goal Coordinate, hook DistanceApiHookFunc) (d *float64, err error) {
	resp, res, err := callDistanceApi(start, goal)
	if err != nil {
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
	if hook != nil {
		hook(start.Id, goal.Id, *resp)
	}
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
	return fmt.Sprintf("%f,%f", c.Long, c.Lat)
}

func MakeDistanceGraph(pairClusters []PairCluster, errorf func(format string, args ...interface{}), hook DistanceApiHookFunc) []DistanceGraph {
	res := []DistanceGraph{}
	for _, pairCluster := range pairClusters {
		dg := DistanceGraph{
			Center:           pairCluster.Center,
			StartCoordinates: pairCluster.Pairs.Coordinates("start"),
			StartIds:         map[string]int{},
			GoalCoordinates:  pairCluster.Pairs.Coordinates("goal"),
			GoalIds:          map[string]int{},
		}
		dg.StartGraph = graph.New(len(dg.StartCoordinates))
		dg.GoalGraph = graph.New(len(dg.GoalCoordinates))

		fillGraph(dg.StartCoordinates, dg.StartGraph, dg.StartIds, errorf, hook)
		fillGraph(dg.GoalCoordinates, dg.GoalGraph, dg.GoalIds, errorf, hook)
		res = append(res, dg)
	}
	return res
}

func fillGraph(currentCs Coordinates, currentGraph *graph.Mutable, ids map[string]int, errorf func(string, ...interface{}), hook DistanceApiHookFunc) {
	for i := 0; i < len(currentCs)-1; i++ {
		for j := i + 1; j < len(currentCs); j++ {
			if d, err := getRoadDistance(currentCs[i], currentCs[j], hook); err != nil {
				if errorf != nil {
					errorf("make distance graph error: %s\n", err)
				}
			} else {
				ids[currentCs[i].Id] = i
				ids[currentCs[j].Id] = j
				currentGraph.AddBothCost(i, j, int64(*d))
			}
		}
	}
}

func AssignDriverToGraphs(graphs []DistanceGraph, drivers Drivers) map[string]*DistanceGraph {
	driverCluster := map[string]*DistanceGraph{}
	for _, d := range drivers {
		driverCluster[d.Id] = nil
	}

	driverPosition := drivers.Coordinates()
	for {
		for driverId, driverCoord := range driverPosition {
			if driverCluster[driverId] == nil {
				distance := math.MaxFloat64
				assignedidx := 0
				for graphidx, graph := range graphs {
					if driverCoord.Distance(graph.Center.Coordinates()) < distance {
						driverCluster[driverId] = &graph
						assignedidx = graphidx
					}
				}
				if len(graphs)-1 <= assignedidx+1 {
					graphs = graphs[:assignedidx]
				} else {
					graphs = append(graphs[:assignedidx], graphs[assignedidx+1:]...)
				}
			}
		}

		nilCounter := 0
		for _, v := range driverCluster {
			if v == nil {
				nilCounter++
			}
		}
		if nilCounter == 0 {
			break
		}
	}
	return driverCluster
}

func FindActions(graphs map[string]*DistanceGraph, req CalculateRequest) CalculateResult {
	res := CalculateResult{Actions: map[string][]DriverAction{}}
	driverPosition := req.Drivers.Coordinates()
	for driverId, currentGraph := range graphs {
		//processing start graph
		firstStart, _ := ClosestCoordinate(driverPosition[driverId], currentGraph.StartCoordinates) //enterance (first) coordinate of start graph
		firstStartIndex := currentGraph.StartIds[firstStart.Id]                                     //graph index of above coordinate
		res.Actions[driverId] = append(res.Actions[driverId], DriverAction{true, *firstStart})      //enterance coordinate is always first action
		var currentCoord Coordinate                                                                 //temp variable
		graph.BFS(currentGraph.StartGraph, firstStartIndex, func(v, w int, c int64) {
			for id, idx := range currentGraph.StartIds {
				if idx == w {
					currentCoord = *currentGraph.StartCoordinates.Search(id)
				}
			}
			res.Actions[driverId] = append(res.Actions[driverId], DriverAction{true, currentCoord})
		})
		finalStart := currentCoord //final assigned temp coordinate is final (exit) coordinate of start graph

		//process goal graph
		firstGoal, _ := ClosestCoordinate(finalStart, currentGraph.GoalCoordinates)            //search enterance coodinate of goal graph
		firstGoalIndex := currentGraph.GoalIds[firstGoal.Id]                                   //graph index of right above coorinate
		res.Actions[driverId] = append(res.Actions[driverId], DriverAction{false, *firstGoal}) //enterance coordinate is always first action
		graph.BFS(currentGraph.GoalGraph, firstGoalIndex, func(v, w int, c int64) {
			for id, idx := range currentGraph.GoalIds {
				if idx == w {
					currentCoord = *currentGraph.GoalCoordinates.Search(id)
				}
			}
			res.Actions[driverId] = append(res.Actions[driverId], DriverAction{false, currentCoord})
		})
	}
	return res
}
