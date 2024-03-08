package response

import (
    "encoding/json"
    "log"
    "net/http"
)

type ApiResponse struct {
    Status string `json:"status"`
    Data   any    `json:"data"`
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, response *ApiResponse) {
    encoder := json.NewEncoder(w)
    w.WriteHeader(statusCode)
    if err := encoder.Encode(*response); err != nil {
        log.Println(err)
    }
}

func NewSuccessResponse(payload any) *ApiResponse {
    return &ApiResponse{
        Status: "success",
        Data:   payload,
    }
}

func NewFailResponse(payload any) *ApiResponse {
    return &ApiResponse{
        Status: "fail",
        Data:   payload,
    }
}

func NewErrorResponse(err error) *ApiResponse {
    return &ApiResponse{
        Status: "error",
        Data:   err.Error(),
    }
}
