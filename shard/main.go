// Copyright (c) 2014 Eric Robert. All rights reserved.

package main

import (
	"bytes"
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
	url       = flag.String("http-address", "", "<addr>:<port> to listen on for HTTP requests")
	reportURL = flag.String("report-url", "", "URL where error reports are posted")
	metricURL = flag.String("metric-url", "", "URL where metrics are posted")
	repeatURL = flag.String("repeat-url", "", "URL where incoming requests are repeated as-is")
	routes    = flag.String("routes", "", "routes configuration")
)

type Route struct {
	Name    string
	Pattern string
	Kind    string
	Sharder json.RawMessage
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
	s := Routes{}
	err = decoder.Decode(&s)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, r := range s {
		d := shard.NewDispatcher(r.Name)

		var s interface{}
		switch {
		case r.Kind == "json":
			s = new(shard.JSONSharder)

		default:
			log.Fatal("route doesn't specify a supported kind of endpoint")
		}

		err = json.Unmarshal(r.Sharder, s)
		if err != nil {
			log.Fatal(err.Error())
		}

		d.Sharder = s.(shard.Sharder)

		if *reportURL != "" {
			d.Reporter = &report.Reporter{
				Name: r.Name,
			}

			d.Reporter.PublishFunc(func(r *report.Report, bodies map[string][]byte) {
				b := new(bytes.Buffer)

				e := json.NewEncoder(b)
				if err := e.Encode(r); err != nil {
					panic(err.Error())
				}

				if err := r.WriteBody(b, bodies); err != nil {
					panic(err.Error())
				}

				_, err = http.Post(*reportURL, "application/json", b)
				if err != nil {
					panic(err.Error())
				}
			})
		}

		if *metricURL != "" {
			d.Monitor = &metric.Monitor{
				Name: r.Name,
			}

			d.Monitor.PublishFunc(func(s *metric.Summary) {
				text, err := json.Marshal(s)
				if err != nil {
					panic(err.Error())
				}

				_, err = http.Post(*metricURL, "application/json", bytes.NewReader(text))
				if err != nil {
					panic(err.Error())
				}
			})
		}

		if *repeatURL != "" {
			d.Repeater = &report.PostRequest{
				URL: *repeatURL,
			}
		}

		d.Start()
		log.Printf("adding route=%s\n", r.Pattern)
		http.Handle(r.Pattern, d)
	}

	log.Printf("starting dispatcher at address=%s\n", *url)
	http.ListenAndServe(*url, nil)
}
