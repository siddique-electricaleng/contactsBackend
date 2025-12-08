package transport

type Envelope struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
