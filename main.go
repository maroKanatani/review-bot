package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"review-bot/structs"
	"review-bot/util"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var token = os.Getenv("TOKEN")
var api = slack.New(token)

func divCheckLine(line string) structs.CheckInfo {
	log.Println(line)
	tokens := strings.Split(line, " ")

	var c structs.CheckInfo

	if len(tokens) < 2 || !(tokens[0] == "[WARN]" || tokens[0] == "[ERROR]") {
		return c
	}

	last := tokens[len(tokens)-1]
	tokens = tokens[:len(tokens)-1]

	c.Level = tokens[0]

	splitted := strings.Split(tokens[1], ":")
	c.FileFullPath = splitted[0]
	// Get file name
	_, f := filepath.Split(splitted[0])
	c.FileName = f
	c.LineNum = splitted[1]
	if !(len(splitted) < 3) {
		c.ColumnNum = splitted[2]
	}

	c.Message = strings.Join(tokens[2:], " ")
	c.CheckType = last
	return c
}

func handle(c echo.Context) error {
	r := c.Request()
	log.Println(r.Header.Get("X-Slack-Retry-Num"))
	if r.Header.Get("X-Slack-Retry-Num") != "" {
		return c.String(200, "OK")
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	fmt.Println(body)
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if e != nil {
		util.ErrLog(e)
		return e
	}
	if eventsAPIEvent.Type == slackevents.URLVerification {
		param := new(structs.Verif)
		if err := json.Unmarshal(buf.Bytes(), &param); err != nil {
			util.ErrLog(err)
			return err
		}
		return c.String(200, param.Challenge)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {

		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			log.Println("---------------------------")
			log.Printf("%+v\n", ev)
			log.Println("---------------------------")
			reqJSON := new(structs.RequestJSON)
			if err := json.Unmarshal(buf.Bytes(), &reqJSON); err != nil {
				util.ErrLog(err)
				return err
			}

			err := onAppMentioned(reqJSON, ev.Channel)
			if err != nil {
				util.ErrLog(err)
				return err
			}

		case *slackevents.MessageEvent:
			log.Println("---------------------------")
			log.Printf("%+v\n", ev)
			log.Println("---------------------------")
			if ev.ChannelType == "im" {
				reqJSON := new(structs.RequestJSON)
				if err := json.Unmarshal(buf.Bytes(), &reqJSON); err != nil {
					util.ErrLog(err)
					return err
				}

				if len(reqJSON.Event.Files) != 0 {
					err := onAppMentioned(reqJSON, ev.Channel)
					if err != nil {
						util.ErrLog(err)
						return err

					}
				}
			}
		}
	}

	return c.String(200, "OK")
}

func CreateTempDirAndFile(dirName string, fName string) (*os.File, error) {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0777)
		if err != nil {
			util.ErrLog(err)
			return nil, err
		}
	}

	_, err = os.Stat(filepath.Join(dirName, fName))
	if os.IsExist(err) {
		util.ErrLog(err)
		return nil, err
	}

	os.Chdir(dirName)
	file, err := os.Create(fName)
	if err != nil {
		util.ErrLog(err)
		return nil, err
	}
	os.Chdir("..")
	return file, nil
}

func onAppMentioned(reqJSON *structs.RequestJSON, channel string) error {
	log.Println("---------------------------")
	log.Printf("%+v\n", reqJSON)
	log.Println("---------------------------")

	if len(reqJSON.Event.Files) == 0 {
		log.Println("No files .")
		_, _, err := api.PostMessage(channel, slack.MsgOptionText("こんにちは！review-botです。\n私にMentionを付けて.javaファイルを送信してください。", false))
		if err != nil {
			util.ErrLog(err)
			return err
		}
	} else {
		dirName := util.NewSecret(32)
		firstTextBlock := slack.NewTextBlockObject(slack.PlainTextType, "お疲れ様です:wave:\nレビュー結果が出ましたのでご確認ください", true, false)
		firstSection := slack.NewSectionBlock(firstTextBlock, nil, nil)
		blocks := []slack.Block{firstSection}

		for _, fStruct := range reqJSON.Event.Files {
			if fStruct.FileType != "java" {
				_, _, err := api.PostMessage(channel, slack.MsgOptionText(":warning:.java以外のファイルが含まれています。\nファイルを見直してから再度.javaファイルを送信してください。", false))
				if err != nil {
					util.ErrLog(err)
					return err
				}
				return nil
			}

			file, err := CreateTempDirAndFile(dirName, fStruct.Name)
			if err != nil {
				return err
			}
			err = api.GetFile(fStruct.URLPrivateDownload, file)
			if err != nil {
				util.ErrLog(err)
				return err
			}
			defer file.Close()

			checkList, err := execCheckStyle(filepath.Join(dirName, fStruct.Name))
			if err != nil {
				util.ErrLog(err)
				return err
			}

			blocks = append(blocks, slack.NewDividerBlock())
			fileNameTextBlock := slack.NewTextBlockObject(slack.MarkdownType, "*"+fStruct.Name+"*", false, false)
			blocks = append(blocks, slack.NewSectionBlock(fileNameTextBlock, nil, nil))
			// lines := strings.Split(string(reviewResult), "\n")
			if len(checkList) == 0 {
				blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, "*指摘は特にありません:ok_hand:*", true, false), nil, nil))
			}

			for _, c := range checkList {
				blocks = append(blocks, c.CreateReviewBlock())
			}
		}

		_, _, err := api.PostMessage(channel, slack.MsgOptionBlocks(blocks...))
		if err != nil {
			util.ErrLog(err)
			return err
		}

	}
	return nil
}

func execCheckStyle(filePath string) ([]structs.CheckInfo, error) {
	const CheckStyleJar = "checkstyle-8.32-all.jar"
	const StyleXML = "mycheck.xml"

	cmd := exec.Command("java", "-jar", CheckStyleJar, "-c", StyleXML, filePath)

	reviewResult, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(reviewResult))
		return nil, err
	}
	lines := strings.Split(string(reviewResult), "\n")

	var checkList []structs.CheckInfo
	for _, line := range lines {
		c := divCheckLine(line)
		if c.Level == "" {
			continue
		}
		checkList = append(checkList, c)
	}
	return checkList, nil
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
	// e.GET("/events-endpoint", handle)
	e.POST("/events-endpoint", handle)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
