package main

import (
	"github.com/labstack/echo/v4"
)

//Coordinate is basic definition of coorinate. It using for directing any senders or stuffs.
type Coordinate struct {
	Id   string  `json:"id"`   //Coordinate's id
	Lat  float64 `json:"lat"`  //Latitude
	Long float64 `json:"long"` //Longitude
}

//Driver expresses current available shipping driver.
type Driver struct {
	Id          string     `json:"id"`           //Driver's name
	Position    Coordinate `json:"position"`     //Driver's current position
	AvailRadius float64    `json:"avail_radius"` //Driver's available moving raduis from current position. Unit is km.
}

//Stuff expresses current requested shipping object which driver ship.
type Stuff struct {
	Id                string     `json:"id"`              //Stuff's name
	SenderName        string     `json:"sender_name"`     //Sender's name
	SenderPosition    Coordinate `json:"sender_position"` //Sender's position
	ReceieverName     string     `json:"recver_name"`     //Receiver's name
	ReceieverPosition Coordinate `json:"recver_position"` //Receiver's position
}

//CalculateRequest is structure for api request.
type CalculateRequest struct {
	Drivers []Driver `json:"drivers"` //Current available drivers data
	Stuffs  []Stuff  `json:"stuffs"`  //Current available stuffs data
}

//DriverAction express every driver's action.
type DriverAction struct {
	IsPickup bool   `json:"is_pickup"` //True if current action is picking stuff up. False if deliver stuff down.
	StuffId  string `json:"stuff_id"`  //Targer stuff's id
}

//CalculateResult is structure for api response.
type CalculateResult struct {
	Actions map[string][]DriverAction `json:"actions"`
}

var (
	validAccessKey []string = make([]string, 0)
	issuedToken    []string = make([]string, 0)
)

func rIssueToken(c echo.Context) error {

}

func validateToken(token string) bool {

}

func rCalculatePath(c echo.Context) error {

}
