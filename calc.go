package main

import (
	"fmt"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"math"
)

type (
	//Coordinate is basic definition of coorinate. It using for directing any senders or stuffs.
	Coordinate struct {
		Id   string  `json:"id" validate:"required"` //Coordinate's id
		Lat  float64 `json:"lat"`                    //Latitude
		Long float64 `json:"long"`                   //Longitude
	}

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

func calculateActions(req CalculateRequest) (*CalculateResult, error) {
	_, err := getKmeanCluster(extractCoordinate("center", &req.Stuffs), len(req.Drivers))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

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
				fmt.Sprintf("c-%d-%d", s.SenderPosition.Id, s.ReceieverPosition.Id),
				(s.SenderPosition.Lat + s.ReceieverPosition.Lat) / 2,
				(s.SenderPosition.Long + s.ReceieverPosition.Long) / 2},
			)
		}
	}
	return res
}

func normalizeCoordinate(c *Coordinate) *Coordinate {
	d := math.Sqrt(math.Pow(c.Lat, 2) + math.Pow(c.Long, 2))
	return &Coordinate{c.Id, c.Lat / d, c.Long / d}
}

func getKmeanCluster(points []*Coordinate, clusterCount int) (clusters.Clusters, error) {
	var d clusters.Observations
	for _, p := range points {
		p = normalizeCoordinate(p)
		d = append(d, clusters.Coordinates{
			p.Lat,
			p.Long,
		})
	}
	km := kmeans.New()
	return km.Partition(d, clusterCount)
}

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
