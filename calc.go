package main

import (
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"strings"
)

func calculateActions(req CalculateRequest) (*CalculateResult, error) {
	_, err := GetKmeanCluster(req.Stuffs.Coordinates("center"), len(req.Drivers))
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
