package api

type response struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
	StatusCode int         `json:"status_code"`
}

type errorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
