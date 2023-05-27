// Package plugindemo a demo plugin.
package statusdonrouters

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

// Config the plugin configuration.
type Config struct {
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	ServerPrefix string `json:"serverPrefix"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
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

type Metric struct {
	Server       string `json:"server"`
	Host         string `json:"host"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	RequestSize  int    `json:"requestSize"`
	ResponseSize int    `json:"responseSize"`
}

func (p *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// create a custom response writer to intercept the response
	crw := &customResponseWriter{ResponseWriter: rw}

	m := &Metric{}

	p.next.ServeHTTP(crw, req)

	// dump the request and get its size
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}

	m.Server = p.config.ServerPrefix
	m.Host = req.Host

	m.Method = req.Method
	m.Path = req.URL.Path
	m.RequestSize = len(dump)
	m.ResponseSize = crw.size

	go p.send(m)
}

func (p *Plugin) send(metric *Metric) {
	v, err := json.Marshal(metric)

	if err != nil {
		// TODO não sei o que fazer aqui
		return
	}

	h := fmt.Sprint(
		p.config.Ip,
		":",
		p.config.Port,
	)

	conn, err := net.Dial("udp", h)

	if err != nil {
		// TODO não sei o que fazer aqui
		return
	}

	defer conn.Close()

	fmt.Fprint(conn, string(v))
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
