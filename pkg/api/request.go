package api

type Request struct {
	Query string      `json:"query"`
	Vars  interface{} `json:"variables"`
}

type responseError struct {
	Message string
}

func (re responseError) Error() string {
	return "DLive API Error: " + re.Message
}

type Response struct {
	Data   map[string]interface{} `json:"data"`
	Errors []responseError        `json:"errors"`
}

type WebSocketRequest struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Payload Request `json:"payload"`
}
