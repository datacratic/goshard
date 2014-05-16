goshard
=======

[![Build Status](https://travis-ci.org/EricRobert/goshard.svg?branch=master)](https://travis-ci.org/EricRobert/goshard)

Installing Go
-------------

Please refer to http://golang.org/doc/install

Source
------

Getting the code:

```
go get github.com/EricRobert/goshard
```

Usage
-----

This tool is designed to provide a simple proxy that splits HTTP requests into partitions (shards) based on the content of each message. For more details see [this](http://en.wikipedia.org/wiki/Shared_nothing_architecture).

To run `shard`, simply supply a binding address and a JSON configuration file.

```
shard --http-address=:12345 --routes=./config.json
```

This JSON configuration file contains an array of routes. Multiple routes can be configured independently by simply providing multiple JSON objects in the top-level array.

Each route has a number of possible parameters:

Field | Description
--- | ---
Name | friendly identifier for this route
Pattern | pattern that will be matched against the URL of the incoming request
Kind | type of sharder e.g. url, http, json (default)
MetricURL | optional URL where metrics are recorded
ReportURL | optional URL where reports are sent
RecordURL | optional URL where incoming requests are sent as content

Then, depending on the type of sharder, a number of parameters are possible.

Field | Description
--- | ---
Shards | number of shards (must be > 0)
URLs | map of shard id to URL
ShardedURL | optional URL that can be built from the shard id e.g. http://server.%d.org
DefaultURL | fallback URL
Field | name of the JSON field used for hashing (json only)
Header | name of the HTTP header used for hashing (http only)

For example, this will install a route that will create 8 partitions based on the `id` field of the JSON message.

```
[
  {
    "Name": "Requests",
    "Kind": "json",
    "Pattern": "/requests",
    "Sharder": {
      "Field": "id",
      "Table": {
        "Shards": 8,
        "ShardedURL": "http://server.%d.org/requests"
      }
    }
  }
]
```

In this other example, the server names are explicitely stated. First, there is a lookup in the URLs table before using the ShardedURL or the DefaultURL to determine where the request will be forwarded. Note that shard indexes are zero-based.

```
[
  {
    "Name": "Requests",
    "Kind": "json",
    "Pattern": "/requests",
    "Sharder": {
      "Field": "id",
      "Table": {
        "Shards": 2,
        "URLs": {
          "0": "http://server-a.org/requests",
          "1": "http://server-b.org/requests"
        }
      }
    }
  }
]
```

Sharders
--------

Type | Description
--- | ---
json | parse the HTTP body as JSON and use specified string or number field as key
http | use the specified HTTP header as key
url | use the last part of the URL as key e.g. /requests/aA1bB2cC3dD4eE5fF6

