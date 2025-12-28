package models

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type TaskCreateResponse struct {
	ID int64 `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Password struct {
	Password string `json:"password"`
}

type JSONToken struct {
	Token string `json:"token"`
}

type EmptyJSONResponse struct {
}
