package main

import (
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"strings"
)

func calculateActions(req CalculateRequest) (*CalculateResult, error) {
	_, err := GetKmeanCluster(extractCoordinate("center", &req.Stuffs), len(req.Drivers))
	if err != nil {
		return nil, err
	}
	return nil, nil
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
				cstart := searchReqCoordById(req, ids[1])
				cgoal := searchReqCoordById(req, ids[2])
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

//Return mock CalculateResult for testing
func mockCalculateActions(req CalculateRequest) (*CalculateResult, error) {
	//Mock result for test
	return &CalculateResult{
		map[string][]DriverAction{
			"0": []DriverAction{
				DriverAction{true, "0"},
				DriverAction{true, "1"},
				DriverAction{false, "1"},
				DriverAction{false, "0"},
			},
			"1": []DriverAction{
				DriverAction{true, "2"},
				DriverAction{false, "2"},
				DriverAction{true, "3"},
				DriverAction{false, "3"},
			},
			"2": []DriverAction{
				DriverAction{true, "4"},
				DriverAction{true, "5"},
				DriverAction{true, "6"},
				DriverAction{false, "6"},
				DriverAction{false, "4"},
				DriverAction{false, "5"},
			},
		},
	}, nil
}
