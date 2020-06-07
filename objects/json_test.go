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
package objects_test

import (
	"bytes"
	"github.com/threeguys/golang-toolkit/objects"
	"github.com/threeguys/golang-toolkit/objects/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

const JsonObjectData = "{\"foo\": \"bar\", \"baz\": { \"red\": \"black\", \"blue\": \"green\" } }"

func TestNewGenericJsonFromReader(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	rdr := &mocks.BadReader{}
	j, err := objects.NewGenericJsonFromReader(rdr)
	assert.Nil(j)
	assert.NotNil(err)
	assert.Equal(mocks.MockErrBadRead, err.Error())

	data := []byte("not-real-json")
	rdr = &mocks.BadReader{ SkipReadError: true, DataRead: data }
	j, err = objects.NewGenericJsonFromReader(rdr)
	assert.Nil(j)
	assert.NotNil(err)

	data = []byte(JsonObjectData)
	rdr = &mocks.BadReader{ SkipReadError: true, DataRead: data }
	j, err = objects.NewGenericJsonFromReader(rdr)
	assert.NotNil(j)
	assert.Nil(err)

	assert.Equal("bar", j.Get("foo"))
}

func TestParseJsonPost_HappyCase(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(JsonObjectData))
	obj, err := objects.ParseJsonPost(r, int64(len(JsonObjectData) + 10))
	assert.NotNil(obj)
	assert.Nil(err)

	r = httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(JsonObjectData))
	obj, err = objects.ParseJsonPost(r, 10)
	assert.Nil(obj)
	assert.NotNil(err)

	r = httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(JsonObjectData))
	obj, err = objects.ParseJsonPost(r, 1024 * 1024)
	assert.NotNil(obj)
	assert.Nil(err)
}

func TestGenericJson_Get(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	obj, err := objects.NewGenericJson([]byte(JsonObjectData))
	assert.Nil(err)

	assert.Equal("bar", obj.Get("foo"))
	assert.Equal("", obj.Get("boh"))
	assert.Equal("bar", obj.GetOrDefault("foo", "should-be-there"))
	assert.Equal("yo", obj.GetOrDefault("boh", "yo"))

	o2 := obj.AsObject("baz")
	assert.NotNil(o2)
	assert.Equal("black", o2.Get("red"))
	assert.Equal("green", o2.Get("blue"))
}

func TestGenericJson_AsList(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	data := []byte("{\"alist\": [{\"a\":\"b\"},{\"c\":\"d\"}]}")

	obj, err := objects.NewGenericJson(data)
	assert.Nil(err)

	assert.Nil(obj.AsList("not-a-name"))

	l := obj.AsList("alist")
	assert.NotNil(l)

}