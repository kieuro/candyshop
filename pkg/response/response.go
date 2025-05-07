package response

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      error  `json:"error"`
}
