package main

import (
	"fmt"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"strings"
)

type (
	//Driver expresses current available shipping driver.
	Driver struct {
		Id          string     `json:"id"`           //Driver's name
		Position    Coordinate `json:"position"`     //Driver's current position
		AvailRadius float64    `json:"avail_radius"` //Driver's available moving raduis from current position. Unit is km.
	}

	//Stuff expresses current requested shipping object which driver ship.
	Stuff struct {
		Id                string     `json:"id"`              //Stuff's name
		SenderName        string     `json:"sender_name"`     //Sender's name
		SenderPosition    Coordinate `json:"sender_position"` //Sender's position
		ReceieverName     string     `json:"recver_name"`     //Receiver's name
		ReceieverPosition Coordinate `json:"recver_position"` //Receiver's position
	}

	//DriverAction express every driver's action.
	DriverAction struct {
		IsPickup bool   `json:"is_pickup"` //True if current action is picking stuff up. False if deliver stuff down.
		StuffId  string `json:"stuff_id"`  //Targer stuff's id
	}
)

/*
	Coordinate is basic definition of coorinate. It using for directing any senders or stuffs.
	It implements (github.com/muesli/clusters).Observation interface for compatible with kmean library
*/
type Coordinate struct {
	Id   string  `json:"id" validate:"required"` //Coordinate's id
	Lat  float64 `json:"lat"`                    //Latitude
	Long float64 `json:"long"`                   //Longitude
}

func (c Coordinate) Coordinates() clusters.Coordinates {
	return clusters.Coordinates([]float64{c.Lat, c.Long})
}

func (c Coordinate) Distance(point clusters.Coordinates) float64 {
	return c.Coordinates().Distance(point)
}

type (
	PairCluster struct {
		Center clusters.Observation
		Pairs  []Pair
	}

	Pair struct {
		Id    string
		Start Coordinate
		Goal  Coordinate
	}
)

func calculateActions(req CalculateRequest) (*CalculateResult, error) {
	_, err := GetKmeanCluster(extractCoordinate("center", &req.Stuffs), len(req.Drivers))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//Extract Coordinates from CalculateRequest with given filtering option
func extractCoordinate(opt string, stuffs *[]Stuff) []*Coordinate {
	res := []*Coordinate{}
	for _, s := range *stuffs {
		switch opt {
		case "send":
			res = append(res, &s.SenderPosition)
		case "recv":
			res = append(res, &s.ReceieverPosition)
		case "center":
			res = append(res, &Coordinate{
				fmt.Sprintf("c-%s-%s", s.SenderPosition.Id, s.ReceieverPosition.Id),
				(s.SenderPosition.Lat + s.ReceieverPosition.Lat) / 2,
				(s.SenderPosition.Long + s.ReceieverPosition.Long) / 2},
			)
		}
	}
	return res
}

//Calculate kmean cluster from given Coordinates
func GetKmeanCluster(points []*Coordinate, clusterCount int) (clusters.Clusters, error) {
	var d clusters.Observations
	for _, p := range points {
		d = append(d, *p)
	}
	km := kmeans.New()
	return km.Partition(d, clusterCount)
}

//Convert "center extracted" coordinate to original coordinates pair with given clusters
func convertCenterToPair(req CalculateRequest, cs clusters.Clusters) []PairCluster {
	res := []PairCluster{}
	for _, c := range cs {
		pc := PairCluster{c.Center, []Pair{}}
		for _, o := range c.Observations {
			oc := o.(Coordinate)
			if strings.Index(oc.Id, "c-") == 0 {
				ids := strings.Split(oc.Id, "-")
				cstart := searchCoordinateById(req, ids[1])
				cgoal := searchCoordinateById(req, ids[2])
				pc.Pairs = append(pc.Pairs, Pair{oc.Id, *cstart, *cgoal})
			}
		}
		res = append(res, pc)
	}
	return res
}

//Search coordinate which has given id from CalcualteRequest
func searchCoordinateById(req CalculateRequest, id string) *Coordinate {
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
