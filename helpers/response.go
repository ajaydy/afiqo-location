package helpers

type (
	Response struct {
		BaseResponse
		Data interface{} `json:"data"`
	}
	BaseResponse struct {
		Errors []string `json:"errors,omitempty"`
	}
)
