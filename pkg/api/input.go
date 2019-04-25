package api

// StreamChatInput represents a message to be sent to a stream's chat
type StreamChatInput struct {
	Message     string `json:"message"`     // The message shown in chat
	RoomRole    string `json:"roomRole"`    // Role of the user sending the chat in the streamer's channel
	Streamer    string `json:"streamer"`    // ID of the streamer
	Subscribing bool   `json:"subscribing"` // Indicates if the user is subscribed to this streamer
}