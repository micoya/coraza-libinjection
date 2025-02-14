//go:build cgo
// +build cgo

// Copyright 2021 Juan Pablo Tosso
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

/*
#cgo LDFLAGS: -L/usr/local/lib -linjection
#include "libinjection.h"
#include "libinjection_sqli.h"
#include "stdlib.h"
*/
import "C"
import (
	"bytes"
	"unsafe"
)

func isSQLi(statement string) (bool, string) {
	var out [8]C.char
	pointer := (*C.char)(unsafe.Pointer(&out[0]))
	cStatement := C.CString(statement)
	defer C.free(unsafe.Pointer(cStatement))
	if found := C.libinjection_sqli(cStatement, C.size_t(len(statement)), pointer); found == 1 {
		output := C.GoBytes(unsafe.Pointer(&out[0]), 8)
		return true, string(output[:bytes.Index(output, []byte{0})])
	}
	return false, ""
}

func isXSS(input string) bool {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	if found := C.libinjection_xss(cInput, C.size_t(len(input))); found == 1 {
		return true
	}
	return false
}
