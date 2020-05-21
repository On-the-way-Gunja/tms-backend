package main

import (
	"encoding/json"
	"fmt"
	"github.com/RyanCarrier/dijkstra"
	"github.com/go-resty/resty/v2"
	"github.com/muesli/clusters"
)

var globalClient *resty.Client

func InitClient() {
	globalClient = resty.New()
}

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

func MakeDistanceGraph(pcs []PairCluster, errorf func(format string, args ...interface{})) ClusterGraph {
	res := map[*clusters.Observation]*dijkstra.Graph{}
	for _, pc := range pcs {
		g := dijkstra.NewGraph()
		cs := extractPairsToCoords(pc.Pairs)

		//Add verticles
		for _, c := range cs {
			g.AddMappedVertex(c.Id)
		}

		//Add arcs
		for i := 0; i < len(cs)-1; i++ {
			for j := i + 1; j < len(cs); j++ {
				if d, err := getRoadDistance(cs[i], cs[j]); err != nil {
					if errorf != nil {
						errorf("make distance graph error: %s\n", err)
					}
				} else {
					g.AddMappedArc(cs[i].Id, cs[j].Id, int64(*d))
				}
			}
		}

		res[&pc.Center] = g
	}
	return res
}

func AssignDriverToGraphs(gs ClusterGraph, drivers []Driver) map[string]*dijkstra.Graph {
	driverCluster := map[string]*dijkstra.Graph{}
	for _, d := range drivers {
		driverCluster[d.Id] = nil
	}

	driverPosition := extractDriversToCoords(drivers)
	for {
		for driverId, driverCoord := range driverPosition {
			if driverCluster[driverId] == nil {
				distance := 9999999999.00
				for gObs, g := range gs {
					if driverCoord.Distance(*gObs.Coordinates()) < distance {
						driverCluster[driverId] = g
					}
				}
			}
		}

		nilCounter := 0
		for k, v := range driverCluster {
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

func FindPath(graphs map[string]*dijkstra.Graph, drivers []Driver) []dijkstra.BestPath {
	for graphName, graph := range graphs {
		graph.Shortest()
	}
}
