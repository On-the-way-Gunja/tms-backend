package main

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

func calculateActions(req CalculateRequest) CalculateResult {
	return CalculateResult{map[string][]DriverAction{"0": []DriverAction{DriverAction{true, "0"}}}} //Mock for test
}
