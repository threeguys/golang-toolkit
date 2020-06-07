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
package clients

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/threeguys/golang-toolkit/servers"
	"io/ioutil"
	"log"
	"net/http"
)

func NewClientMutualTlsConfig(cb *servers.CertificateBundle) (*tls.Config, error) {
	log.Println("Opening root cert", cb.RootPath)
	caCert, err := ioutil.ReadFile(cb.RootPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	log.Println("Opening keypair", cb.CertPath, cb.KeyPath)
	cert, err := tls.LoadX509KeyPair(cb.CertPath, cb.KeyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: caCertPool,
		Certificates: []tls.Certificate{ cert },
	}, nil
}

func NewMutualTlsClient(cb *servers.CertificateBundle) (*http.Client, error) {
	tlsConfig, err := NewClientMutualTlsConfig(cb)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http.Transport{ TLSClientConfig: tlsConfig },
	}, nil
}
