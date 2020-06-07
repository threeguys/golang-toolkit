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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type GenericJson struct {
	fields map[string]interface{}
	value interface{}
}

type GenericJsonList []*GenericJson

func NewGenericJson(data []byte) (*GenericJson, error) {
	var inf interface{}
	err := json.Unmarshal(data, &inf)
	if err != nil {
		return nil, err
	}

	return &GenericJson{
		fields: inf.(map[string]interface{}),
	}, nil
}

func NewGenericJsonFromReader(rdr io.ReadCloser) (*GenericJson, error) {
	defer SafeClose(rdr)

	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}

	return NewGenericJson(data)
}

func NewGenericJsonFromFile(path string) (*GenericJson, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewGenericJson(data)
}

func (gj *GenericJson) GetOrDefault(name, defaultValue string) string {
	if value, ok := gj.fields[name]; ok {
		return value.(string)
	}

	return defaultValue
}

func (gj *GenericJson) Get(name string) string {
	return gj.GetOrDefault(name, "")
}

func (gj *GenericJson) AsObject(name string) *GenericJson {
	if value, ok := gj.fields[name]; ok {
		return &GenericJson{ fields: value.(map[string]interface{}), value: nil }
	}
	return nil
}

func (gj *GenericJson) AsList(name string) GenericJsonList {
	if value, ok := gj.fields[name]; ok {
		items := value.([]interface{})
		jsonItems := make([]*GenericJson, len(items))
		for i, v := range items {
			// TODO handle arrays of primitive types
			jsonItems[i] = &GenericJson{
				fields: v.(map[string]interface{}),
				value:  nil,
			}
		}
		return jsonItems
	}
	return nil
}

func ParseJsonPost(r *http.Request, postLimit int64) (*GenericJson, error) {
	if r.ContentLength > postLimit {
		return nil, errors.New("post content too large")
	}

	return NewGenericJsonFromReader(r.Body)
}

func JsonToMap(data []byte) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
