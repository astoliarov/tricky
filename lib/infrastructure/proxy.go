package infrastructure

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"tricky/lib/usecases"

	"github.com/astoliarov/goproxy"
)

type Proxy struct {
	proxy        *goproxy.ProxyHttpServer
	rulesUseCase *usecases.RulesUseCase
	port         string
	certPath     string
}

func (pr *Proxy) interceptRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	var keys []string
	for key := range r.URL.Query() {
		keys = append(keys, key)
	}

	u, err := pr.rulesUseCase.FindUrl(keys)
	if err != nil {
		return r, nil
	}

	if u == nil {
		return r, nil
	}

	// TODO: think about request modification
	r.URL = u
	r.RequestURI = u.String()
	r.Host = u.Hostname()

	return r, nil
}

func (pr *Proxy) loadCert() (*tls.Certificate, error) {
	caKey, err := ioutil.ReadFile(fmt.Sprintf("%s.key.pem", pr.certPath))
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(fmt.Sprintf("%s.pem", pr.certPath))
	if err != nil {
		return nil, err
	}

	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return nil, err
	}

	return &goproxyCa, nil
}

func (pr *Proxy) Run() error{
	return http.ListenAndServe(pr.port, pr.proxy)
}

func NewProxy(rulesUseCase *usecases.RulesUseCase, port string, certPath string, debug bool) (*Proxy, error) {
	pr := &Proxy{}
	pr.rulesUseCase = rulesUseCase
	pr.certPath = certPath
	pr.port = port

	cert, err := pr.loadCert()
	if err != nil {
		return nil, err
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Ca = cert
	proxy.Verbose = debug
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(pr.interceptRequest)

	pr.proxy = proxy

	return pr, nil
}
