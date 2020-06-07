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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/threeguys/golang-toolkit/objects"
	"io/ioutil"
	"log"
	"net/http"
)

type CertificateBundle struct {
	RootPath string `json:"root"`
	CertPath string `json:"cert"`
	KeyPath string  `json:"key"`
	Domain string   `json:"domain"`
}

//func (cb *CertificateBundle) Load() error {
//
//	rootData, err := ioutil.ReadFile(cb.RootPath)
//	if err != nil {
//		return err
//	}
//
//	_, err = ioutil.ReadFile(cb.CertPath)
//	if err != nil {
//		return err
//	}
//
//	_, err = ioutil.ReadFile(cb.KeyPath)
//	if err != nil {
//		return err
//	}
//
//	cb.rootPem = rootData
//	return nil
//}
//
//func (cb *CertificateBundle) NewCertPool() *x509.CertPool {
//	pool := x509.NewCertPool()
//	pool.AppendCertsFromPEM(cb.rootPem)
//	return pool
//}
//
//func (cb *CertificateBundle) NewServerTlsConfig(pool *x509.CertPool) *tls.Config {
//	config := &tls.Config{
//		ClientCAs: pool,
//		ClientAuth: tls.RequireAndVerifyClientCert,
//		RootCAs: pool,
//	}
//	config.BuildNameToCertificate()
//	return config
//}
//
//func (cb *CertificateBundle) NewTlsServer(bind string, handler http.Handler, config *tls.Config) *http.Server {
//	server := &http.Server{
//		Addr:      bind,
//		TLSConfig: config,
//		Handler:   handler,
//	}
//
//	return server
//}
//
//func (cb *CertificateBundle) ListenAndServe(server *http.Server) error {
//	return server.ListenAndServeTLS(cb.CertPath, cb.KeyPath)
//}

func LoadCertificateBundle(path string) (*CertificateBundle, error) {
	obj, err := objects.NewGenericJsonFromFile(path)
	if err != nil {
		return nil, err
	}

	return &CertificateBundle{
		RootPath: obj.GetOrDefault("root", "ca.pem"),
		CertPath: obj.GetOrDefault("cert", "cert.pem"),
		KeyPath:  obj.GetOrDefault("key", "key.pem"),
		Domain:   obj.GetOrDefault("domain", "(default)"),
	}, nil
}

func NewServerMutualTlsConfig(cb *CertificateBundle) (*tls.Config, error) {
	// Validate the config
	if len(cb.RootPath) == 0 {
		return nil, errors.New("root path must be set")
	} else if len(cb.CertPath) == 0 {
		return nil, errors.New("cert path must be set")
	} else if len(cb.KeyPath) == 0 {
		return nil, errors.New("key path must be set")
	}

	log.Println("Loading root cert:", cb.RootPath)
	// Load the root CA cert
	rootPem, err := ioutil.ReadFile(cb.RootPath)
	if err != nil {
		return nil, err
	}

	// Create a cert pool for the CA
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(rootPem)

	// Create a tls config
	config := &tls.Config{
		ClientCAs: pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		RootCAs: pool,
	}
	config.BuildNameToCertificate()

	return config, nil
}

func ListenAndServerMutualTls(bind string, cb *CertificateBundle, handler http.Handler) error {
	// https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/

	//// Validate the config
	//if len(cb.RootPath) == 0 {
	//	return errors.New("root path must be set")
	//} else if len(cb.CertPath) == 0 {
	//	return errors.New("cert path must be set")
	//} else if len(cb.KeyPath) == 0 {
	//	return errors.New("key path must be set")
	//}
	//
	//// Load the root CA cert
	//rootPem, err := ioutil.ReadFile(cb.RootPath)
	//if err != nil {
	//	return err
	//}
	//
	//// Create a cert pool for the CA
	//pool := x509.NewCertPool()
	//pool.AppendCertsFromPEM(rootPem)
	//
	//// Create a tls config
	//config := &tls.Config{
	//	ClientCAs: pool,
	//	ClientAuth: tls.RequireAndVerifyClientCert,
	//	RootCAs: pool,
	//}
	//config.BuildNameToCertificate()

	config, err := NewServerMutualTlsConfig(cb)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:      bind,
		TLSConfig: config,
		Handler:   handler,
	}

	return server.ListenAndServeTLS(cb.CertPath, cb.KeyPath)
}