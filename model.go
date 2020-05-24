package main

import (
	"fmt"
	"github.com/muesli/clusters"
	"github.com/yourbasic/graph"
	"time"
)

type (
	//Driver expresses current available shipping driver.
	Driver struct {
		Id          string     `json:"id"`           //Driver's name
		Position    Coordinate `json:"position"`     //Driver's current position
		AvailRadius float64    `json:"avail_radius"` //Driver's available moving raduis from current position. Unit is km.
	}
	Drivers []Driver

	//Stuff expresses current requested shipping object which driver ship.
	Stuff struct {
		Id                string     `json:"id"`              //Stuff's name
		SenderName        string     `json:"sender_name"`     //Sender's name
		SenderPosition    Coordinate `json:"sender_position"` //Sender's position
		ReceieverName     string     `json:"recver_name"`     //Receiver's name
		ReceieverPosition Coordinate `json:"recver_position"` //Receiver's position
	}
	Stuffs []Stuff

	//DriverAction express every driver's action.
	DriverAction struct {
		IsPickup   bool       `json:"is_pickup"`  //True if current action is picking stuff up. False if deliver stuff down.
		Coordinate Coordinate `json:"coordinate"` //Action's coordinate
	}

	/*
		Coordinate is basic definition of coorinate. It using for directing any senders or stuffs.
		It implements clusters.Observation interface for compatible with kmean library
	*/
	Coordinate struct {
		Id   string  `json:"id" validate:"required"` //Coordinate's id
		Lat  float64 `json:"lat"`                    //Latitude
		Long float64 `json:"long"`                   //Longitude
	}
	Coordinates []Coordinate

	//KeyRequest is form of issueing token request.
	KeyRequest struct {
		Key string `json:"key" validate:"required"`
	}

	//TokenResponse is form of issueing token response.
	TokenResponse struct {
		Token string `json:"token"`
	}

	//CalculateRequest is structure for api request.
	CalculateRequest struct {
		Drivers Drivers `json:"drivers" validate:"required"` //Current available drivers data
		Stuffs  Stuffs  `json:"stuffs" validate:"required"`  //Current available stuffs data
	}

	//CalculateResult is structure for api response.
	CalculateResult struct {
		Actions          map[string][]DriverAction `json:"actions"` //key is Driver's id
		ActualApiResults map[string]string         `json:"naver_actual_result"`
		EveryApiResults  []PairWithDistance        `json:"naver_every_result"`
	}

	//PairCluster is pair version of clusters.Cluster
	PairCluster struct {
		Center clusters.Observation
		Pairs  Pairs
	}

	//Pair is container of two Coordinates with id.
	Pair struct {
		Id    string
		Start Coordinate
		Goal  Coordinate
	}
	Pairs []Pair

	PairWithDistance struct {
		Pair
		ApiResult *NaverResponse
	}

	//ConfigFormat is definition of required settings
	ConfigFormat struct {
		AccessKey         []string `json:"access_key" validate:"required"`
		NaverClientId     string   `json:"naver_client_id" validate:"required"`
		NaverClientSecret string   `json:"naver_client_secret" validate:"required"`
		EnableTLS         bool     `json:"enable_tls" validate:"required"`
		TLSDomains        []string `json:"tls_domain" validate:"required"`
	}

	//NaverResponse is definition of naver api (https://apidocs.ncloud.com/ko/ai-naver/maps_directions/driving) response
	NaverResponse struct {
		Code            int                    `json:"code"`
		Message         string                 `json:"message"`
		CurrentDateTime string                 `json:"currentDateTime"`
		Route           map[string]interface{} `json:"route"`
	}

	//Token is issued to approved users.
	Token struct {
		Token      string    //Token string
		IssuedTime time.Time //Issued datetime
	}

	DistanceGraph struct {
		Center           clusters.Observation
		StartGraph       *graph.Mutable
		StartCoordinates Coordinates
		StartIds         map[string]int
		GoalGraph        *graph.Mutable
		GoalCoordinates  Coordinates
		GoalIds          map[string]int
	}
	DistanceApiHookFunc func(startId, goalId string, result []byte)
)

func (c Coordinate) Coordinates() clusters.Coordinates {
	return clusters.Coordinates([]float64{c.Lat, c.Long})
}

func (c Coordinate) Distance(point clusters.Coordinates) float64 {
	return c.Coordinates().Distance(point)
}

/*
func (c ClusterGraph) Observations() []*clusters.Observation {
	res := []*clusters.Observation{}
	for k, _ := range c {
		res = append(res, k)
	}
	return res
}
*/

//Extract Coordinates from CalculateRequest with given filtering option
func (s Stuffs) Coordinates(opt string) Coordinates {
	res := []Coordinate{}
	for _, s := range s {
		switch opt {
		case "send":
			res = append(res, s.SenderPosition)
		case "recv":
			res = append(res, s.ReceieverPosition)
		case "center":
			res = append(res, Coordinate{
				fmt.Sprintf("c-%s-%s", s.SenderPosition.Id, s.ReceieverPosition.Id),
				(s.SenderPosition.Lat + s.ReceieverPosition.Lat) / 2,
				(s.SenderPosition.Long + s.ReceieverPosition.Long) / 2},
			)
		case "all":
			res = append(res, s.SenderPosition)
			res = append(res, s.ReceieverPosition)
		}
	}
	return res
}

//Extract all Coordinates from Pairs
func (ps Pairs) Coordinates(opt string) Coordinates {
	res := Coordinates{}
	for _, p := range ps {
		if opt == "all" || opt == "start" {
			if res.Search(p.Start.Id) == nil {
				res = append(res, p.Start)
			}
		}
		if opt == "all" || opt == "goal" {
			if res.Search(p.Goal.Id) == nil {
				res = append(res, p.Goal)
			}
		}
	}
	return res
}

func (drivers Drivers) Coordinates() map[string]Coordinate {
	res := map[string]Coordinate{}
	for _, d := range drivers {
		res[d.Id] = d.Position
	}
	return res
}

//Search coordinate which has given id from Coordinate array
func (cs Coordinates) Search(id string) *Coordinate {
	for _, c := range cs {
		if c.Id == id {
			return &c
		}
	}
	return nil
}
