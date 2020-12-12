package api

type FormMin struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Options map[string]string `json:"options"`
	Class   string            `json:"class"`
}

type DataMin struct {
	Id           string `json:"id"`
	RecordNumber string `json:"record_number"`
	FormId       string `json:"form_id"`
	UserId       string `json:"user_id"`
	CreateTime   string `json:"create_time"`
	AnswerTime   string `json:"answer_time"`
	Direction    string `json:"direction"`
}
