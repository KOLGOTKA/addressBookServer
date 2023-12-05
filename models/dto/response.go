package dto

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

func (resp *Response) GetJson() (byteResp []byte) {
	if resp.Data == nil {
		resp.Data = json.RawMessage(`{}`)
	}
	byteResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error in GetJson") //////////////////////////////////////
	}
	return byteResp
}
