package ghttp

// ResponseBody represent common format response to client
type ResponseBody struct {
	Success    bool   `json:"success"`
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"-"`
	Message    string `json:"message,omitempty"`
	Error      any    `json:"error,omitempty"`
	Paging     any    `json:"paging,omitempty"`
}
