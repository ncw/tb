# Token Bucket (tb)

This package provides a generic implementation of the "Token bucket"
algorithm where the handling of non-conformity is left to the user.
[Wikipedia page](http://en.wikipedia.org/wiki/Token_bucket)
![Image](http://sardes.inrialpes.fr/~krakowia/MW-Book/Chapters/QoS/Chapters/QoS/Figs/bucket.gif)


## Install
```shell
$ go get github.com/tsenart/tb
```

## Usage
```go
// A an example of an echo server which throttles each remote ip to 10
// connections per second.
package main

import (
  "github.com/tsenart/tb"
  "io"
  "log"
  "net"
)

func main() {
	ln, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatal(err)
	}

	echo := func(conn net.Conn) {
		defer conn.Close()

		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}
		// Throttle to 10 connection per second from the same host
		// Handle non-conformity by dropping the connection
		if out := tb.Throttle(host, 1, 10); out < 1 {
			log.Printf("Throttled %s", host)
			return
		}
		log.Printf("Echoing payload from %s:%s", host, port)
		io.Copy(conn, conn)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go echo(conn)
	}
}
```

## Licence
```
The MIT License (MIT)

Copyright (c) 2014 Tomás Senart

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
