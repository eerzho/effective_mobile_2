package response

type Error struct {
	Error  string   `json:"error"`
	Errors []string `json:"errors,omitempty"`
}

type Success struct {
	Data interface{} `json:"data,omitempty"`
}
