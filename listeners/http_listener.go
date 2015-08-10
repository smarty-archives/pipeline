package listeners

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HTTPListener struct {
	certificatePEM string
	inner          http.Server
}

func NewHTTPListener(listenAddress string, handler http.Handler) *HTTPListener {
	return &HTTPListener{
		inner: http.Server{
			Addr:           listenAddress,
			Handler:        handler,
			ReadTimeout:    time.Second * 15,
			WriteTimeout:   time.Second * 15,
			MaxHeaderBytes: 1024 * 2,
			ErrorLog:       log.New(ioutil.Discard, "", 0),
		},
	}
}
func (this *HTTPListener) WithTLS(certificatePEM string, tlsConfig *tls.Config) {
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			SessionTicketsDisabled:   true,
			CipherSuites: []uint16{
				tls.TLS_FALLBACK_SCSV,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	this.certificatePEM = certificatePEM
	this.inner.TLSConfig = tlsConfig
}

func (this *HTTPListener) Listen() {
	log.Printf("[INFO] Listening for web traffic on %s.\n", this.inner.Addr)
	if err := this.listen(); err != nil {
		log.Fatal("[ERROR] Unable to listen to web traffic: ", err)
	}
}
func (this *HTTPListener) listen() error {
	if len(this.certificatePEM) == 0 {
		return this.inner.ListenAndServe()
	}

	return this.inner.ListenAndServeTLS(this.certificatePEM, this.certificatePEM)
}
