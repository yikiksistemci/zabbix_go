package zabbix

import (
	"fmt"
)

type APIError struct {
	Code int `json:"code"`

	Message string `json:"message"`

	Data string `json:"data"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}
