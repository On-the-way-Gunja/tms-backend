package main

import (
	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/on-the-way-gunja/tms-backend/docs"
	"github.com/swaggo/echo-swagger"
	"github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

//build : ~/go/bin/swag init && go build && sudo ./tms*

//Struct validator
type Validator struct {
	validator *validator.Validate
}

//Validate function
func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var Config *ConfigFormat

// @title OTW Prototype API
// @version 1.0
// @contact.email kimdictor@gmail.com
// @description On the Way's api prototype for demonstraning technology.

// @license.name exclusive-closed license

func main() {
	//Initiate echo
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	//Set logging
	f := new(prefixed.TextFormatter)
	f.FullTimestamp = true
	elogrus.Attach(e).Logger.Formatter = f

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{"API-TOKEN"},
	}))

	//Read config
	if c, err := ReadConfig("config.json", e.Validator.Validate); err != nil {
		e.Logger.Fatal(err)
	} else {
		Config = c
	}

	if Config.EnableTLS {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(Config.TLSDomains...)
		e.AutoTLSManager.Cache = autocert.DirCache(".cache")
	}

	InitMapClient()

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.POST("/token", rIssueToken)
	e.POST("/path", rCalculatePath)
	if Config.EnableTLS {
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	} else {
		e.Logger.Fatal(e.Start(":80"))
	}
}
