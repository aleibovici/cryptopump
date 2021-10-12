## HOW TO COMPILE

To compile cryptopump you need Go environment set.

On Ubuntu 18.04 you can install Go as follow:

Download the Go language binary archive:
```
$ wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
```

Extract it:
```
$ sudo tar -xvf go1.16.4.linux-amd64.tar.gz
```

and copy it:
```
$ sudo mv go /usr/local
```

Setup Go Environment:
```
$ export GOROOT=/usr/local/go
$ export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```
Verify Go is running with
```
$ go version
```
and check the output for
```
go version go1.16.4 linux/amd64
```

Now go to cryptopump directory and compile it with
```
$ go build .
```

An executable should be present in cryptopump directory.
