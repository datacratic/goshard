// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"bytes"
	"fmt"
	"github.com/datacratic/goreports"
	"net/http"
	"time"
)

type Sharder interface {
	Shard(req *report.Request) (string, int, error)
}

type Dispatcher struct {
	report.Endpoint
	Sharder
	Client http.Client
}

type DispatcherMetrics struct {
	Request        bool
	ReceivedFailed bool
	Invalid        bool
	PostFailed     bool
	SendFailed     bool
	ResponseFailed bool
	Failed         bool
	ReadDuration   time.Duration
	PostDuration   time.Duration
	FullDuration   time.Duration
	Shard          string
	StatusCode     string
}

func (d *Dispatcher) routeMessage(w http.ResponseWriter, r *http.Request) (metrics DispatcherMetrics) {
	metrics.Request = true

	t0 := time.Now()

	req, err := report.NewRequest(r)
	if err != nil {
		d.ReportErrorWithRequest(err, req)
		metrics.ReceivedFailed = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.RecordRequest(req)

	url, k, err := d.Sharder.Shard(req)
	if err != nil {
		d.ReportErrorWithRequest(err, req)
		metrics.Invalid = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.Shard = fmt.Sprintf("%d", k)

	t1 := time.Now()
	metrics.ReadDuration = t1.Sub(t0)

	q, err := http.NewRequest(r.Method, url, bytes.NewReader(req.Content))
	if err != nil {
		d.ReportErrorWithRequest(err, req)
		metrics.PostFailed = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q.Header = r.Header

	p, err := d.Client.Do(q)
	if err != nil {
		d.ReportErrorWithRequest(err, req)
		metrics.SendFailed = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rep, err := report.NewResponse(p)
	if err != nil {
		d.ReportErrorWithRequestAndResponse(err, req, rep)
		metrics.ResponseFailed = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.StatusCode = fmt.Sprintf("%d", p.StatusCode)

	if p.StatusCode != http.StatusOK && p.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("http response=%d", p.StatusCode)
		d.ReportErrorWithRequestAndResponse(err, req, rep)
		metrics.Failed = true
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.PostDuration = time.Now().Sub(t1)

	if p.StatusCode == http.StatusOK {
		for k, v := range p.Header {
			w.Header().Set(k, v[0])
			for i := 1; i != len(v); i++ {
				w.Header().Add(k, v[i])
			}
		}

		w.Write(rep.Content)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	metrics.FullDuration = time.Now().Sub(t0)
	return
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := d.routeMessage(w, r)
	d.RecordMetrics(&result)
}

func NewDispatcher(name string) *Dispatcher {
	d := new(Dispatcher)
	d.Name = name
	return d
}
