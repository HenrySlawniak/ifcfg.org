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
	"fmt"
	"github.com/dgryski/go-identicon"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"net/http"
	"strings"
	"time"
)

func main() {
	cLog := console.New()
	cLog.SetTimestampFormat(time.RFC3339)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Info("Starting ifcfg.org")

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("ifcfg.org", "v4.ifcfg.org", "v6.ifcfg.org"),
		Cache:      autocert.DirCache("certs"),
	}

	httpSrv := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(httpRedirectHandler),
	}

	go httpSrv.ListenAndServe()

	rootSrv := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   http.HandlerFunc(pageHandler),
	}

	http2.ConfigureServer(rootSrv, &http2.Server{})

	rootSrv.ListenAndServeTLS("", "")
}

func httpRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(split[:len(split)-1], ":")
	// This is bad, and I feel bad
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	log.Infof("Handling %s from url %s\n\tUA: %s\n", ip, r.URL.String(), r.Header.Get("User-Agent"))
	if strings.Contains(r.URL.String(), "favicon.ico") {
		w.Header().Set("Content-Type", "image/png")
		w.Write(generateIco([]byte(ip)))
		return
	}

	if strings.Contains(r.Header.Get("User-Agent"), "curl") || r.Header.Get("Accepts") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip + "\n"))
		return
	}
	w.Write([]byte(ip))
}

func generateIco(dat []byte) []byte {
	dat = append(dat, []byte(fmt.Sprintf("%x", time.Now().UnixNano()))[:]...)
	return identicon.New7x7([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}).Render(dat)
}
