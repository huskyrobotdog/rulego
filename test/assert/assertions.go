/*
 * Copyright 2023 The RuleGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// CallerInfo This function is inspired by:
// https://github.com/stretchr/testify/blob/master/assert/assertions.go
func CallerInfo() []string {
	var pc uintptr
	var ok bool
	var file string
	var line int
	var name string

	callers := []string{}
	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case, see #180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		parts := strings.Split(file, "/")
		if len(parts) > 1 {
			filename := parts[len(parts)-1]
			dir := parts[len(parts)-2]
			if (dir != "assert" && dir != "mock" && dir != "require") || filename == "mock_test.go" {
				callers = append(callers, fmt.Sprintf("%s:%d", file, line)+"\n")
			}
		}

	}

	return callers
}

func Equal(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v\n Error Trace:   %s", a, b, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}

func EqualCleanString(t *testing.T, a, b string) {
	re := regexp.MustCompile(`[\s\n\t]+`)
	cleanStra := re.ReplaceAllString(a, "")
	cleanStrb := re.ReplaceAllString(b, "")
	if cleanStra != cleanStrb {
		t.Errorf("%v ！= %v\n Error Trace:   %s", a, b, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func NotEqual(t *testing.T, a, b interface{}) {
	if reflect.DeepEqual(a, b) {
		t.Errorf("%v == %v\n Error Trace:   %s", a, b, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}

func True(t *testing.T, value bool) {
	if !value {
		t.Errorf("%v should be true\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func False(t *testing.T, value bool) {
	if value {
		t.Errorf("%v should be false\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func NotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Errorf("%v should be not nil\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
func Nil(t *testing.T, value interface{}) {
	if value != nil {
		t.Errorf("%v should be nil\n Error Trace:   %s", value, strings.Join(CallerInfo(), "\n\t\t\t"))
	}
}
