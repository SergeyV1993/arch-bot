package models

type Update struct {
	ID      int
	Message *IncomingMessage
}

type IncomingMessage struct {
	Text     string
	Username string
	ChatID   int64

	Location *Location
}

type Location struct {
	Longitude string
	Latitude  string
}
