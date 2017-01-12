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
	"encoding/json"
	"encoding/xml"
	"flag"
	"github.com/HenrySlawniak/go-identicon"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"image/color"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

const version = "1.1.0"

var domains = flag.String("domain", "ifcfg.org,v4.ifcfg.org,v6.ifcfg.org", "A comma-seperaated list of domains to get a certificate for.")
var client = &http.Client{}

func main() {
	flag.Parse()
	cLog := console.New()
	cLog.SetTimestampFormat(time.RFC3339)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Info("Starting ifcfg.org")

	domainList := strings.Split(*domains, ",")
	for i, d := range domainList {
		domainList[i] = strings.TrimSpace(d)
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domainList...),
		Cache:      autocert.DirCache("certs"),
	}

	httpSrv := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(httpRedirectHandler),
	}

	go httpSrv.ListenAndServe()

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/detail", detailHandler)
	mux.HandleFunc("/favicon.ico", icoHandler)

	rootSrv := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   mux,
	}

	http2.ConfigureServer(rootSrv, &http2.Server{})

	rootSrv.ListenAndServeTLS("", "")
}

func httpRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

// GetIP returns the remote ip of the request, by stripping off the port from the RemoteAddr
func GetIP(r *http.Request) string {
	split := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(split[:len(split)-1], ":")
	// This is bad, and I feel bad
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	return ip
}

func icoHandler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	w.Header().Set("Content-Type", "image/png")
	w.Write(generateIco([]byte(ip)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	w.Header().Set("Server", "ifcfg.org v"+version)

	if strings.Contains(r.Header.Get("User-Agent"), "curl") || r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip + "\n"))
		return
	}
	w.Write([]byte(ip))
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	w.Header().Set("Server", "ifcfg.org v"+version)
	detail, err := ARINLookup(ip)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Encountered an error: " + err.Error()))
		debug.PrintStack()
		log.Error(err)
		return
	}
	j, err := json.MarshalIndent(detail, "", "  ")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Encountered an error: " + err.Error()))
		debug.PrintStack()
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(j)
}

func generateIco(dat []byte) []byte {
	return identicon.New7x7([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}).RenderWithBG(dat, color.NRGBA{0x0, 0x0, 0x0, 0x0})
}

func ARINLookup(ip string) (*ARINroot, error) {
	url := "https://whois.arin.net/rest/ip/" + ip
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/xml")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dat := ARINroot{}
	err = xml.Unmarshal(bod, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}
