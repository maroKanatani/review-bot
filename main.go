package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var token = os.Getenv("TOKEN")
var api = slack.New(token)

type URLVeri struct {
	token     string `json:"token"`
	challenge string `json:"challenge"`
	_type     string `json:"type"`
}

func handle(c echo.Context) error {
	log.Println("Hello")
	r := c.Request()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	fmt.Println(body)
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if e != nil {
		errLog(e)
		return e
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var param URLVeri
		if err := c.Bind(param); err != nil {
			errLog(err)
			return err
		}
		return c.String(200, param.challenge)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		}
	}

	return c.String(200, "OK")
}

func main() {
	// fmt.Println("[INFO] Server listening")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	// http.ListenAndServe(":"+port, nil)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/events-endpoint", handle)
	e.POST("/events-endpoint", handle)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

func errLog(err error) {
	log.Printf("%+v\n", errors.WithStack(err))
}
