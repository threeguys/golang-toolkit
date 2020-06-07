//
//  Copyright 2020 Ray Cole
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
package objects

import (
    "errors"
    "fmt"
    "reflect"
    "testing"
)

type Assertions struct {
    t *testing.T
}

func NewTestAssertions(t *testing.T) *Assertions {
    return &Assertions{ t: t }
}

func (a *Assertions) Fail(msg string) {
    a.t.Helper()
    if a.t == nil {
        panic(errors.New(msg))
    } else {
        a.t.Error(msg)
        a.t.Fail()
    }
}

func (a *Assertions) Equal(x1 interface{}, x2 interface{}) {
    a.t.Helper()
    if !reflect.DeepEqual(x1, x2) {
        a.Fail(fmt.Sprintf("Expected values to equal: %+v != %+v", x1, x2))
    }
}

func (a *Assertions) Nil(x interface{}) {
    a.t.Helper()
    if x != nil && !reflect.ValueOf(x).IsNil(){
        a.Fail(fmt.Sprintf("Expected nil but was %+v", x))
    }
}

func (a *Assertions) NotNil(x interface{}) {
    a.t.Helper()
    if x == nil || reflect.ValueOf(x).IsNil() {
        a.Fail("Expected not nil but was nil")
    }
}

func (a *Assertions) True(answer bool) {
    a.t.Helper()
    a.Equal(answer, true)
}

func (a *Assertions) False(answer bool) {
    a.t.Helper()
    a.Equal(answer, false)
}
