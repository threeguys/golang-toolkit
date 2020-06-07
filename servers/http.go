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
package servers

import (
    "log"
    "net/http"
)

const JSON = "application/json"

const HeaderContentType = "Content-Type"
const HeaderContentEncoding = "Content-Encoding"
const HeaderContentLength = "Content-Length"
const HeaderUserAgent = "User-Agent"

const ContentTypeHtml = "text/html"
const ContentTypeJson = "application/json"
const ContentTypeFormEncoded = "application/x-www-form-urlencoded"


type HttpLogger struct {
    prefix string
    delegate http.Handler
}

func NewHttpLogger(prefix string, delegate http.Handler) *HttpLogger {
    return &HttpLogger {
        prefix: prefix,
        delegate: delegate,
    }
}

func (h *HttpLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("(%s) HTTP %s %s [%s]", h.prefix, r.Method, r.URL.Path, r.Host)
    if r.ContentLength > 0 {
        log.Printf("    BODY[%s] %d", r.Header.Get("Content-Type"), r.ContentLength)
    }
    h.delegate.ServeHTTP(w, r)
}
