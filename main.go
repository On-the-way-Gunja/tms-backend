package main

import (
	"encoding/json"
	elogrus "github.com/dictor/echologrus"
	"github.com/labstack/echo/v4"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
)

func main() {
	e := echo.New()

	//Set logging
	f := new(prefixed.TextFormatter)
	f.FullTimestamp = true
	elogrus.Attach(e).Logger.Formatter = f

	e.Logger.Fatal(e.Start(":80"))
}

func setAccessKey(path string) error {
	if d, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return json.Unmarshal(d, validAccessKey)
	}
}
