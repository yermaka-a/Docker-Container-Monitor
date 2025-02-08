package models

type Container struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	TimeMs   string `json:"timeMs"`
	PingDate string `json:"pingDate"`
}
