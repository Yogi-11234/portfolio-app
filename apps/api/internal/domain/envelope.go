package domain

// Meta holds pagination metadata
type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}

// SuccessResponse defines the standard success envelope
type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
}

// ErrorDetail defines the error object within ErrorResponse
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// ErrorResponse defines the standard error envelope
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// NewSuccessResponse creates a success envelope
func NewSuccessResponse(data interface{}, meta *Meta) SuccessResponse {
	return SuccessResponse{
		Data: data,
		Meta: meta,
	}
}

// NewErrorResponse creates an error envelope
func NewErrorResponse(code, message string, fields map[string]string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}
}
