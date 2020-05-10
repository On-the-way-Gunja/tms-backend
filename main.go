package main

import (
	"encoding/json"
	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/on-the-way-gunja/prototype/docs"
	"github.com/swaggo/echo-swagger"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
)

//build : /home/ubuntu/go/bin/swag init && go build && sudo ./proto*

//Struct validator
type Validator struct {
	validator *validator.Validate
}

//Validate function
func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// @title OTW Prototype API
// @version 1.0
// @contact.email kimdictor@gmail.com
// @description On the Way's api prototype for demonstraning technology.

// @license.name exclusive-closed license

func main() {
	e := echo.New()

	//Set logging
	f := new(prefixed.TextFormatter)
	f.FullTimestamp = true
	elogrus.Attach(e).Logger.Formatter = f

	e.Validator = &Validator{validator: validator.New()}

	if err := setAccessKey("keys.txt"); err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.POST("/token", rIssueToken)
	e.POST("/path", rCalculatePath)
	e.Logger.Fatal(e.Start(":80"))
}

func setAccessKey(path string) error {
	if d, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return json.Unmarshal(d, &validAccessKey)
	}
}
