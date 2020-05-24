package main

import (
	"fmt"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"math"
	"strings"
)

func calculateActions(req CalculateRequest) (*CalculateResult, error) {
	kmean, err := GetKmeanCluster(req.Stuffs.Coordinates("center"), len(req.Drivers))
	if err != nil {
		return nil, err
	}
	pairs := convertCenterToPairs(req, kmean)

	apiResults := map[string]string{}
	hook := func(startId, goalId string, result []byte) {
		apiResults[fmt.Sprintf("%s-%s", startId, goalId)] = string(result)
	}
	errors := ""
	graphs := MakeDistanceGraph(pairs, func(format string, args ...interface{}) {
		errors += fmt.Sprintf(format, args)
	}, hook)
	if errors != "" {
		return nil, fmt.Errorf("MakeDistanceGraph : %s", errors)
	}

	graphWithDrivers := AssignDriverToGraphs(graphs, req.Drivers)
	res := FindActions(graphWithDrivers, req)
	res.ActualApiResults = apiResults
	res.EveryApiResults = calculateAllDistance(req)
	return &res, nil
}

func calculateAllDistance(req CalculateRequest) []PairWithDistance {
	coords := req.Stuffs.Coordinates("all")
	res := []PairWithDistance{}
	for i := 0; i < len(coords)-1; i++ {
		for j := i + 1; j < len(coords); j++ {
			_, r, _ := callDistanceApi(coords[i], coords[j])
			res = append(res, PairWithDistance{
				Pair: Pair{
					Id:    coords[i].Id + "-" + coords[j].Id,
					Start: coords[i],
					Goal:  coords[j],
				},
				ApiResult: r,
			})
		}
	}
	return res
}

//Calculate kmean cluster from given Coordinates
func GetKmeanCluster(points []Coordinate, clusterCount int) (clusters.Clusters, error) {
	var d clusters.Observations
	for _, p := range points {
		d = append(d, p)
	}
	km := kmeans.New()
	return km.Partition(d, clusterCount)
}

//Convert "center extracted" coordinate to original coordinates pair with given clusters
func convertCenterToPairs(req CalculateRequest, cs clusters.Clusters) []PairCluster {
	res := []PairCluster{}
	for _, c := range cs {
		pc := PairCluster{c.Center, []Pair{}}
		for _, o := range c.Observations {
			oc := o.(Coordinate)
			if strings.Index(oc.Id, "c-") == 0 {
				ids := strings.Split(oc.Id, "-")
				cstart := searchCoordFromReq(req, ids[1])
				cgoal := searchCoordFromReq(req, ids[2])
				pc.Pairs = append(pc.Pairs, Pair{oc.Id, *cstart, *cgoal})
			}
		}
		res = append(res, pc)
	}
	return res
}

//Search coordinate which has given id from CalcualteRequest
func searchCoordFromReq(req CalculateRequest, id string) *Coordinate {
	for _, s := range req.Stuffs {
		if s.SenderPosition.Id == id {
			return &s.SenderPosition
		}
		if s.ReceieverPosition.Id == id {
			return &s.ReceieverPosition
		}
	}
	for _, d := range req.Drivers {
		if d.Position.Id == id {
			return &d.Position
		}
	}
	return nil
}

func ClosestCoordinate(src Coordinate, dest Coordinates) (res *Coordinate, dist float64) {
	dist = math.MaxFloat64
	for _, c := range dest {
		if d := src.Distance(c.Coordinates()); d <= dist {
			dist = d
			res = &c
		}
	}
	return
}
