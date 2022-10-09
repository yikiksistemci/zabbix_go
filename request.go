package zabbix

import "sync/atomic"

type Request struct {
	JSONRPCVersion string      `json:"jsonrpc"`
	Method         string      `json:"method"`
	Params         interface{} `json:"params"`
	RequestID      uint64      `json:"id"`
	AuthToken      string      `json:"auth,omitempty"`
}

var requestID uint64

func NewRequest(method string, params interface{}) *Request {
	if params == nil {
		params = map[string]string{}
	}
	return &Request{
		JSONRPCVersion: "2.0",
		Method:         method,
		Params:         params,
		RequestID:      atomic.AddUint64(&requestID, 1),
		AuthToken:      "",
	}
}
