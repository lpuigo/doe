package extracter

type FormConfig struct {
	FormId   int    `json:"form_id"`
	FormName string `json:"form_name"`
}

type Config struct {
	ProgressFile    string       `json:"progress_file"`
	ProgressXLSFile string       `json:"progress_xls_file"`
	Forms           []FormConfig `js:"forms"`
}
