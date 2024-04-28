package response

type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Data interface{} `json:"data,omitempty"`
}
