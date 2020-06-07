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
package system

import (
    "errors"
    "fmt"
    "os"
    "strconv"
    "time"
)

func EnvOrDefault(name string, defaultValue string) string {
    if value, found := os.LookupEnv(name); found {
        return value
    }

    return defaultValue
}

func EnvOrDefaultDuration(name string, defaultValue time.Duration) time.Duration {
    return time.Duration(EnvOrDefaultInt(name, int(defaultValue)))
}

func EnvOrDefaultInt(name string, defaultValue int) int {
    if value, found := os.LookupEnv(name); found {
        intValue, err := strconv.Atoi(value)
        if err == nil {
            return intValue
        } else {
            panic(err)
        }
    }
    return defaultValue
}

func EnvRequired(name string) (string, error) {
    if value, found := os.LookupEnv(name); found {
        return value, nil
    }

    return "", errors.New(fmt.Sprintf("Environment %s was missing", name))
}

func EnvMapRequired(names []string) (map[string]string, error) {
    values := make(map[string]string)
    for _, n := range names {
        v, err := EnvRequired(n)
        if err != nil {
            return nil, err
        }
        values[n] = v
    }
    return values, nil
}
