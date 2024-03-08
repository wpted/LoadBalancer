package response

import (
    "encoding/json"
    "log"
    "net/http"
)

// ApiResponse is the general struct of all server responses.
type ApiResponse struct {
    Status string `json:"status"`
    Data   any    `json:"data"`
}

// WriteJsonResponse encodes the json response to w.
func WriteJsonResponse(w http.ResponseWriter, statusCode int, response *ApiResponse) {
    encoder := json.NewEncoder(w)
    w.WriteHeader(statusCode)
    if err := encoder.Encode(*response); err != nil {
        log.Println(err)
    }
}

// NewSuccessResponse constructs a response when an API call is successful.
func NewSuccessResponse(payload any) *ApiResponse {
    return &ApiResponse{
        Status: "success",
        Data:   payload,
    }
}

// NewFailResponse constructs a response when an API call is rejected due to invalid data or call conditions.
func NewFailResponse(payload any) *ApiResponse {
    return &ApiResponse{
        Status: "fail",
        Data:   payload,
    }
}

// NewErrorResponse constructs a response when failure occurs due to an error on the server.
func NewErrorResponse(err error) *ApiResponse {
    return &ApiResponse{
        Status: "error",
        Data:   err.Error(),
    }
}
