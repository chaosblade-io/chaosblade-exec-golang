package response

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Code    Code        `json:"code"`
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type Code int32

const (
	OK                        = 200
	ExperimentNotFound        = 1404
	ExperimentInCircuit       = 1405
	ExperimentMatcherNotFound = 1406
	ExperimentNotMatched      = 1407
	ExperimentLimited         = 1408
	RequestHandlerNotFound    = 1501
	EncodeError               = 1512
	DecodeError               = 1513
	IllegalParameters         = 1602
)

func ReturnOK(result interface{}) Response {
	return Response{Code: OK, Success: true, Result: result}
}

func ReturnFail(code Code, err string) Response {
	return Response{Code: code, Success: false, Error: err}
}

func (response Response) Print() string {
	bytes, err := json.Marshal(&response)
	if err != nil {
		return fmt.Sprintf("marshall response err, %s; code: %d", err.Error(), response.Code)
	}
	return string(bytes)
}

// Decode return the response that wraps the content
func Decode(content string, defaultValue Response) Response {
	var resp Response
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		defaultValue = ReturnFail(DecodeError, fmt.Sprintf("unmarshal %s err: %v", content, err))
		//logrus.Warningf("decode %s err, return default value, %s", content, defaultValue.Print())
		return defaultValue
	}
	return resp
}
