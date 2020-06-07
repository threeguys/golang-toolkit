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
	"errors"
	"github.com/threeguys/golang-toolkit/objects"
	"github.com/threeguys/golang-toolkit/objects/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSafeClose(t *testing.T) {
	assert := objects.NewTestAssertions(t)

	mc := &mocks.MockCloser{ Err: nil }
	objects.SafeClose(mc)
	assert.Equal(1, mc.Closes)

	mc = &mocks.MockCloser{Err: errors.New("this is a  test")}
	objects.SafeClose(mc)
	assert.Equal(1, mc.Closes)
}

func TestSuccessful(t *testing.T) {
	assert := objects.NewTestAssertions(t)

	w := httptest.NewRecorder()
	assert.True(objects.Successful(nil, w))
	assert.Equal(http.StatusOK, w.Result().StatusCode)

	w = httptest.NewRecorder()
	assert.False(objects.Successful(errors.New("test error"), w))
	assert.Equal(http.StatusInternalServerError, w.Result().StatusCode)
}
