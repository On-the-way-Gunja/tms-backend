package main

import (
	"encoding/json"
	elogrus "github.com/dictor/echologrus"
	"github.com/labstack/echo/v4"
	_ "github.com/on-the-way-gunja/prototype/docs"
	"github.com/swaggo/echo-swagger"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
)

// @title OWT Prototype API
// @version 1.0
// @contact.email kimdictor@gmail.com
// @license.name exclusive-closed

func main() {
	e := echo.New()

	//Set logging
	f := new(prefixed.TextFormatter)
	f.FullTimestamp = true
	elogrus.Attach(e).Logger.Formatter = f

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
