package skynetstomp

type SkynetStompRequest struct {
	FirstName string
	LastName  string
}

type SkynetStompResponse struct {
	Errors []string
	Status string
}
