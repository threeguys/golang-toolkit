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
package mocks

import (
	"errors"
	"io"
)

const MockErrBadRead = "error: bad read"
const MockErrBadClose = "error: bad close"

type MockCloser struct {
	Err error
	Closes int
}

func (mc *MockCloser) Close() error {
	mc.Closes++
	return mc.Err
}

type BadReader struct {
	SkipReadError bool
	DataRead []byte
	doneRead bool
}

func (b *BadReader) Read(p []byte) (n int, err error) {
	if b.SkipReadError {
		if b.DataRead != nil && !b.doneRead {
			b.doneRead = true
			copy(p, b.DataRead)
			return len(b.DataRead), nil
		} else if b.doneRead {
			return 0, io.EOF
		} else {
			b.doneRead = true
			return 0, nil
		}
	}
	return 0, errors.New(MockErrBadRead)
}

func (b *BadReader) Close() error {
	return errors.New(MockErrBadClose)
}
