package api

import (
	"strings"
)

type Field struct {
	Parent   string            `json:"parent"`
	Value    string            `json:"value"`
	Type     string            `json:"type"`
	Subtype  string            `json:"subtype"`
	Hidden   string            `json:"hidden"`
	Multiple string            `json:"multiple"`
	Boolean  bool              `json:"boolean"`
	Time     map[string]string `json:"time"`
}

func (f *Field) String() string {
	res := strings.Builder{}
	res.WriteString(f.Type + "\n")
	if f.Type == "photo" {
		for item, timestamp := range f.Time {
			res.WriteString("\t\tImgUrl " + item + " : " + timestamp + "\n")
		}
	} else {
		res.WriteString("\t\t" + f.Value + "\n")
	}
	return res.String()
}

type FormulaireData struct {
	ID               string `json:"id"`
	RecordNumber     string `json:"record_number"`
	FormID           string `json:"form_id"`
	UserID           string `json:"user_id"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
	UpdateUserID     string `json:"update_user_id"`
	UpdateAnswerTime string `json:"update_answer_time"`
	//StartTime        interface{} `json:"start_time"`
	//EndTime          interface{} `json:"end_time"`
	//Direction        interface{} `json:"direction"`
	//RecipientID      interface{} `json:"recipient_id"`
	History        string            `json:"history"`
	FormUniqueID   string            `json:"form_unique_id"`
	OriginAnswer   string            `json:"origin_answer"`
	AnswerTime     string            `json:"answer_time"`
	UserName       string            `json:"user_name"`
	LastName       string            `json:"last_name"`
	FirstName      string            `json:"first_name"`
	Phone          string            `json:"phone"`
	Email          string            `json:"email"`
	Login          string            `json:"login"`
	UpdateUserName string            `json:"update_user_name"`
	RecipientName  string            `json:"recipient_name"`
	Fields         map[string]*Field `js:"fields"`
}

func (fd *FormulaireData) String() string {
	res := strings.Builder{}
	res.WriteString("Form " + fd.RecordNumber + " (" + fd.UpdateAnswerTime + ") by " + fd.UpdateUserName + "\n")
	for fieldName, field := range fd.Fields {
		if field.Hidden == "true" {
			continue
		}
		if field.Type == "section" {
			continue
		}
		if field.Value == "" {
			continue
		}
		res.WriteString("\t" + fieldName + ": " + field.String())
	}
	return res.String()
}
