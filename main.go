// Copyright (c) 2016 Henry Slawniak <https://henry.computer/>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

func main() {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("ifcfg.org", "v4.ifcfg.org", "v6.ifcfg.org"),
		Cache:      autocert.DirCache("certs"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	rootSrv := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	v4Mux := http.NewServeMux()
	v4Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("v4 Hello world"))
	})

	v6Mux := http.NewServeMux()
	v6Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("v6 Hello world"))
	})

	v4Srv := &http.Server{
		Addr:      "v4.ifcfg.org:https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   v4Mux,
	}

	v6Srv := &http.Server{
		Addr:      "v6.ifcfg.org:https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   v6Mux,
	}

	go rootSrv.ListenAndServeTLS("", "")
	go v4Srv.ListenAndServeTLS("", "")
	v6Srv.ListenAndServeTLS("", "")
}
