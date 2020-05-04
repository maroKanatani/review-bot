package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

type Verif struct {
	Token     string `json:"token" form:"token" query:"token"`
	Challenge string `json:"challenge" form:"challenge" query:"challenge"`
	Type      string `json:"type" form:"type" query:"type"`
}

type RequestJSON struct {
	Token    string `json:"token" form:"token" query:"token"`
	TeamID   string `json:"team_id" form:"team_id" query:"team_id"`
	ApiAppID string `json:"api_app_id" form:"api_app_id" query:"api_app_id"`
	Event    Event  `json:"event" form:"event" query:"event"`
}

type Event struct {
	Type  string `json:"type" form:"type" query:"type"`
	Text  string `json:"text" form:"text" query:"text"`
	Files []File `json:"files" form:"files" query:"files"`
}

type File struct {
	ID                 string `json:"id" form:"id" query:"id"`
	Name               string `json:"name" form:"name" query:"name"`
	FileType           string `json:"filetype" form:"filetype" query:"filetype"`
	User               string `json:"user" form:"user" query:"user"`
	URLPrivateDownload string `json:"url_private_download" form:"url_private_download" query:"url_private_download"`
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
		param := new(Verif)
		if err := json.Unmarshal(buf.Bytes(), &param); err != nil {
			errLog(err)
			return err
		}
		return c.String(200, param.Challenge)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {

		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			reqJSON := new(RequestJSON)
			if err := json.Unmarshal(buf.Bytes(), &reqJSON); err != nil {
				errLog(err)
				return err
			}
			log.Println("---------------------------")
			log.Printf("%+v\n", reqJSON)
			log.Println("---------------------------")
			log.Println("---------------------------")
			log.Printf("%+v\n", ev)
			log.Println("---------------------------")

			if len(reqJSON.Event.Files) == 0 {
				log.Println("No files .")
			} else {
				dirName := newSecret(32)
				err := os.Mkdir(dirName, 0777)
				if err != nil {
					errLog(err)
					return err
				}
				os.Chdir(dirName)
				file, err := os.Create(reqJSON.Event.Files[0].Name)
				if err != nil {
					errLog(err)
					return err
				}
				err = api.GetFile(reqJSON.Event.Files[0].URLPrivateDownload, file)
				if err != nil {
					errLog(err)
					return err
				}
				file.Close()
				// os.Chdir("..")
				// err = os.RemoveAll(dirName)
				// if err != nil {
				// 	errLog(err)
				// 	return err
				// }
			}

			_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				errLog(err)
				return err
			}
		}
	}

	return c.String(200, "OK")
}

func main() {

	fmt.Println("[INFO] Server listening")
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

func newSecret(length int) string {
	runes := make([]byte, length)

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(255))
		runes[i] = byte(num.Int64())
	}

	return base64.RawStdEncoding.EncodeToString(runes)
}
