// +build !prod

package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"cirello.io/svc/jwt"
	"golang.org/x/crypto/acme/autocert"
)

var acceptableTargets = map[string]struct{}{}

func services() {
	log.Println("bootstrapping services")

	m := &autocert.Manager{
		Cache:  autocert.DirCache("./httpd-services.secrets"),
		Prompt: autocert.AcceptTOS,
		Email:  "user@example.com",
		// HostPolicy: autocert.HostWhitelist(),
	}
	log.Println("starting svc:80")
	go func() {
		log.Println("svc:80", http.ListenAndServe(
			servicesBindIP+":http",
			m.HTTPHandler(http.HandlerFunc(nil))))
	}()

	certBytes, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatalln("Unable to read ca.pem", err)
	}

	clientCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalln("unable to load system CA pool", err)
	}
	if ok := clientCAs.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	s := &http.Server{
		Addr: servicesBindIP + ":https",
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
			ClientAuth:     tls.VerifyClientCertIfGiven,
			ClientCAs:      clientCAs,
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cert := detectedClientCertificate(r); cert != nil {
				token, err := jwt.CreateFromCert(r.Host, certBytes, cert)
				if err == nil {
					r.Header.Set("Authorization", "bearer "+token)
				}
			} else if cookie, err := r.Cookie(gatewayTokenCookie); err != nil || cookie.Value == "" {
				handleSSOLogin(r.Host, certBytes, w, r)
				return
			} else {
				r.Header.Set("Authorization", "bearer "+cookie.Value)
			}

			// Add here handlers that need protection.
			http.NotFound(w, r)
		}),
	}
	log.Println("starting svc:443")
	log.Println("svc:443", s.ListenAndServeTLS("", ""))
}
