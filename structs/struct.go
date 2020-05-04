package structs

import "fmt"

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

type CheckInfo struct {
	Level     string
	FileName  string
	LineNum   string
	ColumnNum string
	Message   string
	CheckType string
}

func (c CheckInfo) ShowInfo() {
	fmt.Println("############################")
	fmt.Println("【File】" + c.FileName)
	fmt.Println("【Message】" + c.Message)
	fmt.Println("【Line Num】" + c.LineNum)
	fmt.Println("【Column】" + c.ColumnNum)
	fmt.Println("【Type】" + c.CheckType)
}
