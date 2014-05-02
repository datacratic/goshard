goshard
=======

Installing Go
-------------

Please refer to http://golang.org/doc/install

For example, this is what I did:

```
wget https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz
tar -xzf go1.2.1.linux-amd64.tar.gz
```

And adding this to `~/.profile`

```
export GOROOT=$HOME/go
export GOPATH=$HOME/code/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Then, `source ~/.profile` and you're ready to GO.

Shard
-----

Getting the code:

```
go get github.com/EricRobert/goshard
```

Build:

```
go install github.com/EricRobert/goshard/...
```

Running:

...
