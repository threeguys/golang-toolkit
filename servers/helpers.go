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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/threeguys/golang-toolkit/objects"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseProducer func() (url string, resp *http.Response, err error)


func ExecuteUrlCallback(req *http.Request, cb func(*http.Response)) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error notifying of database failover", err)
		return
	}

	defer objects.SafeClose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Println("Non-200 response from async url", req.URL, resp.StatusCode, resp.Status)
	}

	if cb != nil {
		cb(resp)

	} else if resp.Body != nil {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body from the application")
		}
	}
}


func WriteHttpError(w http.ResponseWriter, message string) error {
	content := []byte(fmt.Sprintf("<html><body><h2>%s</h2></body></html>", message))
	w.Header().Add(HeaderContentType, ContentTypeHtml)
	w.Header().Add(HeaderContentLength, fmt.Sprintf("%d", len(content)))
	w.WriteHeader(http.StatusInternalServerError)

	amt, err := w.Write(content)
	if err != nil {
		return err
	} else if amt != len(content) {
		return errors.New(fmt.Sprintf("Expected to write %d during error response but wrote %d", len(content), amt))
	}
	return nil
}

func SafeWriteHttpError(w http.ResponseWriter, message string) {
	err := WriteHttpError(w, message)
	if err != nil {
		log.Println(err)
	}
}

func WriteHttpContent(w http.ResponseWriter, contentType string, content []byte) error {
	w.Header().Add(HeaderContentType, contentType)
	w.Header().Add(HeaderContentLength, fmt.Sprintf("%d", len(content)))
	w.WriteHeader(http.StatusOK)

	amt, err := w.Write(content)
	if err != nil {
		return err
	} else if amt != len(content) {
		return errors.New(fmt.Sprintf("Expected to write %d during response but wrote %d", len(content), amt))
	}

	return nil
}

func SafeWriteHttpObject(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		SafeWriteHttpError(w, err.Error())
		log.Println(err)
		return
	}

	SafeWriteHttpContent(w, JSON, data)
}

func SafeWriteHttpContent(w http.ResponseWriter, contentType string, content []byte) {
	err := WriteHttpContent(w, contentType, content)
	if err != nil {
		log.Println(err)
	}
}

func NewProducer(url string) ResponseProducer {
	return NewClientProducer(url, http.DefaultClient)
}

func NewClientProducer(url string, client *http.Client) ResponseProducer {
	return func() (string, *http.Response, error) {
		resp, err := client.Get(url)
		return url, resp, err
	}
}

func GetAsJson(url string) (*objects.GenericJson, error) {
	return GetProducerAsJson(NewProducer(url))
}

func GetProducerAsJson(getter ResponseProducer) (*objects.GenericJson, error) {
	url, resp, err := getter()
	if err != nil {
		return nil, err
	}

	defer objects.SafeClose(resp.Body)

	if hdr := resp.Header.Get(HeaderContentType); len(hdr) > 0 {
		if hdr != ContentTypeJson {
			return nil, errors.New(fmt.Sprintf("invalid Content-Type %s for url %s", hdr, url))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("no Content-Type specified in the response to %s", url))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return objects.NewGenericJson(data)
}

func HttpRespond(code int, w http.ResponseWriter) {
	w.Header().Set(HeaderContentLength, "0")
	w.WriteHeader(code)
}

func InternalError(w http.ResponseWriter) {
	HttpRespond(http.StatusInternalServerError, w)
}

func Ok(w http.ResponseWriter) {
	HttpRespond(http.StatusOK, w)
}

func BadRequest(w http.ResponseWriter) {
	HttpRespond(http.StatusBadRequest, w)
}