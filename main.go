package main

import (
	"flag"
	"github.com/RollingBalls/rollingballs-server/webserver"
)

func main() {
	listen := flag.String("listen", ":8080", "Listen address and port")
	flag.Parse()

	webserver.Run(*listen)
}
