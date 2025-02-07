package db

type Container struct {
	ID       int    `json:"id"`
	Ip       string `json:"ip"`
	TimeMs   string `json:"timeMs"`
	PingDate string `json:"pingDate"`
}
