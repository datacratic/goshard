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

This tool is designed to provide a simple proxy that splits HTTP requests into partitions (shards) based on the content of each message. For more details on the utility of this process, see this article about [sharding](http://en.wikipedia.org/wiki/Shared_nothing_architecture).

To run `shard`, simply supply a binding address and a JSON configuration file.

```
shard --http-address=:12345 --routes=./config.json
```

This JSON configuration file contains an array of routes. Multiple routes can be configured independently by simply providing multiple JSON objects in the top-level array. Each route has a number of possible parameters:

Field | Description
--- | ---
Name | friendly identifier for this route
Pattern | pattern that will be matched against the URL of the incoming request
Kind | type of sharder e.g. json (default)
MetricURL | optional url where metrics are recorded
ReportURL | optional url where reports are sent
RecordURL | optional url where incoming requests are sent as content

Then, depending on the type of sharder, a number of parameters are possible.

Field | Description
--- | ---
Shards | number of shards (must be > 0)
URLs | map of shard id to url
ShardedURL | optional url that can be built from the shard id e.g. http://server.%d.org
DefaultURL | fallback url
Field | name of the JSON field used for hashing (json only)

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

