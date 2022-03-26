package handler

type validationError struct {
	Attribute string `json:"attribute"`
	Message   string `json:"message"`
}

type errorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details []validationError `json:"details"`
}
