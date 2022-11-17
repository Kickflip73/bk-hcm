/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

// Package errf defines common error.
package errf

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"hcm/pkg/kit"
	"hcm/pkg/logs"
)

// ErrorF defines an error with error code and message.
type ErrorF struct {
	// Code is hcm errCode
	Code int32 `json:"code"`
	// Message is error detail
	Message string `json:"message"`
}

// Error implement the golang's basic error interface
func (e *ErrorF) Error() string {
	if e == nil || e.Code == OK {
		return "nil"
	}

	// return with a json format string error, so that the upper service
	// can use Wrap to decode it.
	return fmt.Sprintf(`{"code": %d, "message": "%s"}`, e.Code, e.Message)
}

// Format the ErrorF error to a string format.
func (e *ErrorF) Format() string {
	if e == nil || e.Code == OK {
		return ""
	}

	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// AssignResp used only to assign the values of the Code and Message
// fields of ErrorF to the Code and Message fields of the response.
// Node: resp must be a *struct.
func (e ErrorF) AssignResp(kit *kit.Kit, resp interface{}) {
	if reflect.ValueOf(resp).Type().Kind() != reflect.Ptr {
		logs.ErrorDepthf(1, "response is not pointer, rid: %s", kit.Rid)
		return
	}

	if _, ok := reflect.TypeOf(resp).Elem().FieldByName("Code"); !ok {
		logs.ErrorDepthf(1, "response have not 'Code' field, rid: %s", kit.Rid)
		return
	}

	if _, ok := reflect.TypeOf(resp).Elem().FieldByName("Message"); !ok {
		logs.ErrorDepthf(1, "response have not 'Message' field, rid: %s", kit.Rid)
		return
	}

	valueOf := reflect.ValueOf(resp).Elem()

	code := valueOf.FieldByName("Code")
	code.SetInt(int64(e.Code))

	msg := valueOf.FieldByName("Message")
	msg.SetString(e.Message)
}

// New an error with error code and message.
func New(code int32, message string) error {
	return &ErrorF{Code: code, Message: message}
}

// Newf create an error with error code and formatted message.
func Newf(code int32, format string, args ...interface{}) error {
	return &ErrorF{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Error try to convert the error to ErrorF if possible.
// it is used by the RPC client to wrap the response error response
// by the RPC server to the ErrorF, user can use this ErrorF to test
// if an error is returned or not, if yes, then use the ErrorF to
// response with error code and message.
func Error(err error) *ErrorF {
	if err == nil {
		return nil
	}

	ef, ok := err.(*ErrorF)
	if ok {
		return ef
	}

	s := err.Error()

	// test if the error is a json error,
	// if not, then this is an error without error code.
	if !strings.HasPrefix(s, "{") {
		return &ErrorF{
			Code:    Unknown,
			Message: s,
		}
	}

	// this is a standard error format, then decode it directly.
	ef = new(ErrorF)
	if err := json.Unmarshal([]byte(s), ef); err != nil {
		return &ErrorF{
			Code:    Unknown,
			Message: s,
		}
	}

	return ef
}
