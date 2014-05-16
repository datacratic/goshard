// Copyright (c) 2014 Eric Robert. All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"github.com/EricRobert/gometrics"
	"github.com/EricRobert/goreports"
	"github.com/EricRobert/goshard"
	"log"
	"net/http"
	"os"
)

var (
	url    = flag.String("http-address", "", "<addr>:<port> to listen on for HTTP requests")
	routes = flag.String("routes", "", "routes configuration")
)

type Route struct {
	Name      string
	Pattern   string
	Kind      string
	MetricURL string
	ReportURL string
	RecordURL string
	Sharder   json.RawMessage
}

type Routes []Route

func main() {
	flag.Parse()

	if *url == "" {
		log.Fatal("--http-address is required")
	}

	if *routes == "" {
		log.Fatal("--routes is required")
	}

	file, err := os.Open(*routes)
	if err != nil {
		log.Fatal(err.Error())
	}

	decoder := json.NewDecoder(file)
	items := Routes{}
	err = decoder.Decode(&items)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, r := range items {
		d := shard.NewDispatcher(r.Name)

		var s interface{}
		switch {
		case r.Kind == "" || r.Kind == "json":
			s = new(shard.JSONSharder)

		case r.Kind == "http":
			s = new(shard.HTTPSharder)

		case r.Kind == "url":
			s = new(shard.URLSharder)

		case r.Kind == "regexp":
			s = new(shard.RegexpSharder)

		default:
			log.Fatal("route doesn't specify a supported kind of endpoint")
		}

		err = json.Unmarshal(r.Sharder, s)
		if err != nil {
			log.Fatal(err.Error())
		}

		d.Sharder = s.(shard.Sharder)

		if r.MetricURL != "" {
			d.Monitor = metric.NewJSONMonitor(r.Name, r.MetricURL)
		}

		if r.ReportURL != "" {
			d.Reporter = report.NewJSONReporter(r.Name, r.ReportURL)
		}

		if r.RecordURL != "" {
			d.Recorder = &report.PostRequest{
				URL: r.RecordURL,
			}
		}

		if "" == r.Pattern {
			r.Pattern = "/"
		}

		d.Start()
		log.Printf("adding route=%s\n", r.Pattern)
		http.Handle(r.Pattern, d)
	}

	log.Printf("starting dispatcher at address=%s\n", *url)
	http.ListenAndServe(*url, nil)
}
