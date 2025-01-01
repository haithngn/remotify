package values

import "encoding/json"

const genericError = `
{
	"error": {
		"code": 426,
		"message": "Something went wrong"
	}
}`

type failedResponse struct {
	Error cause `json:"error" binding:"required"`
	Code  int   `json:"code" binding:"required"`
}

type cause struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type succeedResponse struct {
	Code    int         `json:"code" binding:"required"`
	Data    interface{} `json:"data" binding:"required"`
	Message string      `json:"message" binding:"required"`
}

// Failed generates a JSON response for a failed operation.
//
// It takes in an error and a code as parameters. The error is used to generate
// the error message, and the code is used to specify the error code.
//
// The function returns a string representing the JSON response.
func Failed(theError error, code int) string {
	b, err := json.Marshal(failedResponse{
		Error: cause{
			Code:    code,
			Message: theError.Error(),
		},
	})

	if err != nil {
		return genericError
	}

	return string(b)
}

// Succeed generates a JSON string representation of a success response with the provided data.
//
// Parameters:
// - data: The data to be included in the response. It can be any valid Go data type.
//
// Returns:
// - string: The JSON string representation of the success response.
func Succeed(data interface{}) string {
	b, err := json.Marshal(succeedResponse{
		Code:    200,
		Data:    data,
		Message: "Success",
	})

	if err != nil {
		return genericError
	}

	return string(b)
}
