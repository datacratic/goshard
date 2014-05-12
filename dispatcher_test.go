// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEndpoint(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	t.Log("starting fake endpoint")
	h := httptest.NewServer(http.HandlerFunc(handler))

	e := NewEndpoint("Test")

	e.Sharder = JsonSharder{
		Field: "id",
		Table: Table{
			Shards:     1,
			DefaultUrl: h.URL,
		},
	}

	e.Start()

	t.Log("starting dispatcher endpoint")
	g := httptest.NewServer(e)

	text := `{"id":"1234567890"}`

	for i := 0; i != 10; i++ {
		r, err := http.Post(g.URL, "application/json", strings.NewReader(text))
		if err != nil {
			t.Fail()
		}

		t.Log(fmt.Sprintf("%+v", r))
	}

	g.Close()
	h.Close()
}
