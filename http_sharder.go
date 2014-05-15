// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"github.com/EricRobert/goreports"
)

type HTTPSharder struct {
	Table  Table
	Header string
}

func (h *HTTPSharder) Shard(req *report.Request) (url string, k int, err error) {
	key := req.Request.Header.Get(h.Header)
	if key == "" {
		err = fmt.Errorf("HTML header '%s' is missing", h.Header)
		return
	}

	url, k = h.Table.GetURL([]byte(key))
	return
}
