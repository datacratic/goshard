// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"github.com/EricRobert/goreports"
	"path"
)

type URLSharder struct {
	Table Table
}

func (h *URLSharder) Shard(req *report.Request) (url string, k int, err error) {
	p := req.Request.URL.Path

	key := path.Base(p)
	if p == "." || p == "/" {
		err = fmt.Errorf("URL '%s' does not ends with an id", p)
		return
	}

	url, k = h.Table.GetURL([]byte(key))
	return
}
