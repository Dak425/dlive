package api

type request struct {
	Query string                 `json:"query"`
	Vars  map[string]interface{} `json:"variables"`
}

type responseError struct {
	Message string
}

func (re responseError) Error() string {
	return "GraphQL API Error: " + re.Message
}

type response struct {
	Data   interface{}
	Errors []responseError
}

type webSocketRequest struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Payload request `json:"payload"`
}