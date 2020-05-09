package main

import (
	elogrus "github.com/dictor/echologrus"
	"github.com/labstack/echo/v4"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
	"net/http"
)

func main() {
	e := echo.New()

	//Set logging
	f := new(prefixed.TextFormatter)
	f.FullTimestamp = true
	elogrus.Attach(e).Logger.Formatter = f

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":80"))
}

func setAccessKey(path string) error {
	if d, err := ioutil.ReadFile; d != nil {
		return err
	} else {
		return json.Unmarshal(string(d), validAccessKey)
	}
}
