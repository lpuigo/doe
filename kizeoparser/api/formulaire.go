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
	res.WriteString("Form " + fd.ID + " by " + fd.UpdateUserName + "\n")
	res.WriteString("CreateTime       : " + fd.CreateTime + "\n")
	res.WriteString("UpdateTime       : " + fd.UpdateTime + "\n")
	res.WriteString("UpdateAnswerTime : " + fd.UpdateAnswerTime + "\n")
	res.WriteString("AnswerTime       : " + fd.AnswerTime + "\n")
	res.WriteString("History          : " + fd.History + "\n")
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

type SearchData struct {
	ID               string `json:"_id"`
	RecordNumber     string `json:"_record_number"`
	FormID           string `json:"_form_id"`
	UserID           string `json:"_user_id"`
	UpdateUserID     string `json:"_update_user_id"`
	CreateTime       string `json:"_create_time"`
	UpdateTime       string `json:"_update_time"`
	AnswerTime       string `json:"_answer_time"`
	UpdateAnswerTime string `json:"_update_answer_time"`
	History          string `json:"_history"`
	FormUniqueID     string `json:"_form_unique_id"`
	UserName         string `json:"_user_name"`
	UpdateUserName   string `json:"_update_user_name"`
	SummaryTitle     string `json:"_summary_title"`
	SummarySubtitle  string `json:"_summary_subtitle"`
	Comment          string `json:"commentaire"`
	Geoloc           struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
	} `json:"_geoloc"`
}

type SearchResult struct {
	RecordsTotal    int           `json:"recordsTotal"`
	RecordsFiltered int           `json:"recordsFiltered"`
	Data            []*SearchData `json:"data"`
}
