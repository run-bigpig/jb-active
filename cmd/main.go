package main

import (
	"flag"
	"github.com/run-bigpig/jb-active/internal/bootstrap"
	"github.com/valyala/fasthttp"
	"log"
)

var addr = flag.String("addr", ":10800", "TCP address to listen to")

func main() {
	flag.Parse()
	r := bootstrap.Run()
	log.Printf("服务已启动,监听于 %s \n", *addr)
	err := fasthttp.ListenAndServe(*addr, r.Handler)
	if err != nil {
		panic(err)
	}
}
