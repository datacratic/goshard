// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonSharder struct {
	Table Table
	Field string
}

func (h JsonSharder) Shard(content []byte) (url string, k int, err error) {
	item := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewBuffer(content))
	err = decoder.Decode(&item)
	if err != nil {
		return
	}

	key, ok := item[h.Field]
	if !ok {
		err = fmt.Errorf("message body is missing a json object with field '%s'", h.Field)
		return
	}

	switch t := key.(type) {
	case float64:
		url, k = h.Table.GetUrl([]byte(fmt.Sprintf("%d", t)))
	case string:
		url, k = h.Table.GetUrl([]byte(t))
	default:
		err = fmt.Errorf("field %s must be a string or a number", h.Field)
	}

	return
}
