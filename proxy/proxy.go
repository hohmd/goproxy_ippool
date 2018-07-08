package proxy

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/*
Provider - proxy provider. Fabric for proxy clients
*/
type Provider struct {
	proxyServerAddress string // HTTP proxy server address
}

/*
Client - proxy client that connects to Proxy service and gets new IP for proxy
Otherwise - inits default http.Client
*/
type Client struct {
	key    string
	client *http.Client
}

// New - Provider constructor
func New(proxyServerAddr string) *Provider {
	return &Provider{
		proxyServerAddress: proxyServerAddr,
	}
}

// NewClient - Client constructor
func (p *Provider) NewClient(key string) *Client {
	c := Client{
		key:    key,
		client: p.obtain(key),
	}
	return &c
}

// obtain - getting new IP address
func (p *Provider) obtain(key string) (client *http.Client) {
	var err error
	client = &http.Client{}
	proxyAddr, err := getProxyAddress(p.proxyServerAddress, key)
	if err != nil {
		log.Println("[ERROR] Getting proxy IP: ", err)
	}
	// Setting Proxy
	var proxyURL *url.URL
	if proxyURL, err = url.Parse("http://" + proxyAddr); err == nil {
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Disable HTTP/2.
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}
	}
	if err != nil {
		fmt.Println("Error parsing proxy IP: ", err)
	}

	return
}

// Do - HTTP request doer
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func getProxyAddress(addr, key string) (ip string, err error) {
	proxyServerAddr := addr + "/" + key
	if !strings.Contains(proxyServerAddr, "http") {
		proxyServerAddr = "http://" + proxyServerAddr
	}
	res, err := http.Get(proxyServerAddr)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	ip = string(b)
	return
}
