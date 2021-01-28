package api

import (
	"fmt"
	"strings"
)

type FormMin struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Options map[string]string `json:"options"`
	Class   string            `json:"class"`
}

func (fm *FormMin) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Id: %8s  \tName:%s\n", fm.Id, fm.Name))
	sb.WriteString(fmt.Sprintf("Class    : %s\n", fm.Class))
	sb.WriteString(fmt.Sprintf("Options  :\n"))
	for attr, value := range fm.Options {
		sb.WriteString(fmt.Sprintf("\t%20s : %s\n", attr, value))
	}
	return sb.String()
}

type DataMin struct {
	Id           string `json:"id"`
	RecordNumber string `json:"record_number"`
	FormId       string `json:"form_id"`
	UserId       string `json:"user_id"`
	CreateTime   string `json:"create_time"`
	AnswerTime   string `json:"answer_time"` //  Date d'enregistrement
	Direction    string `json:"direction"`
}

func (dm *DataMin) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Id           : %s\n", dm.Id))
	sb.WriteString(fmt.Sprintf("RecordNumber : %s\n", dm.RecordNumber))
	sb.WriteString(fmt.Sprintf("FormId       : %s\n", dm.FormId))
	sb.WriteString(fmt.Sprintf("UserId       : %s\n", dm.UserId))
	sb.WriteString(fmt.Sprintf("CreateTime   : %s\n", dm.CreateTime))
	sb.WriteString(fmt.Sprintf("AnswerTime   : %s\n", dm.AnswerTime))
	sb.WriteString(fmt.Sprintf("Direction    : %s\n", dm.Direction))
	return sb.String()
}
