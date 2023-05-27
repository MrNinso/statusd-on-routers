// Package plugindemo a demo plugin.
package statusdonrouters

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Config the plugin configuration.
type Config struct {
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	ServerPrefix string `json:"serverPrefix"`
	RotuerPrefix string `json:"rotuerPrefix"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Ip:           "localhost",
		Port:         "8000",
		ServerPrefix: "traefik",
		RotuerPrefix: "test",
	}
}

// Demo a Demo plugin.
type Plugin struct {
	next   http.Handler
	name   string
	config *Config
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Plugin{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (p *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// create a custom response writer to intercept the response
	crw := &customResponseWriter{ResponseWriter: rw}
	// dump the request and get its size
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}

	reqSize := len(dump)

	p.next.ServeHTTP(crw, req)

	resSize := crw.size

	fmt.Printf("Request size: %d bytes\n", reqSize)
	fmt.Printf("Response size: %d bytes\n", resSize)
}

type customResponseWriter struct {
	http.ResponseWriter
	size int
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	n, err := crw.ResponseWriter.Write(b)
	crw.size += n
	return n, err
}
