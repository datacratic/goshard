// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"bytes"
	"fmt"
	"github.com/EricRobert/goer"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type Sharder interface {
	Shard(content []byte) (string, int, error)
}

type Endpoint struct {
	*service.Endpoint
	Sharder
	Client http.Client
}

type EndpointMetrics struct {
	Request        bool
	ReceivedFailed bool
	Invalid        bool
	PostFailed     bool
	ResponseFailed bool
	Failed         bool
	ReadDuration   time.Duration
	PostDuration   time.Duration
	FullDuration   time.Duration
	Shard          string
	StatusCode     string
}

func (e *Endpoint) routeMessage(w http.ResponseWriter, r *http.Request) (metrics EndpointMetrics) {
	metrics.Request = true

	t0 := time.Now()

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		metrics.ReceivedFailed = true
		dump, _ := httputil.DumpRequest(r, true)
		e.Report(err, dump)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.Body.Close()
	e.Repeat(request, r)

	url, k, err := e.Sharder.Shard(request)
	if err != nil {
		metrics.Invalid = true
		e.Report(err, request)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.Shard = fmt.Sprintf("%d", k)

	t1 := time.Now()
	metrics.ReadDuration = t1.Sub(t0)

	s, err := e.Client.Post(url, r.Header.Get("Content-Type"), bytes.NewReader(request))
	if err != nil {
		metrics.PostFailed = true
		e.Report(err, request)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := ioutil.ReadAll(s.Body)
	if err != nil {
		metrics.ResponseFailed = true
		dump, _ := httputil.DumpResponse(s, true)
		e.Report(err, dump)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.StatusCode = fmt.Sprintf("%d", s.StatusCode)
	s.Body.Close()

	if s.StatusCode != http.StatusOK && s.StatusCode != http.StatusNoContent {
		metrics.Failed = true
		e.Report(fmt.Errorf("http response=%d", s.StatusCode), response)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.PostDuration = time.Now().Sub(t1)

	w.Header().Set("Content-Type", s.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))
	w.Write(response)

	metrics.FullDuration = time.Now().Sub(t0)
	return
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := e.routeMessage(w, r)
	e.Record(&result)
}

func NewEndpoint(name string) *Endpoint {
	e := Endpoint{
		Endpoint: service.NewEndpoint(name),
	}

	return &e
}
