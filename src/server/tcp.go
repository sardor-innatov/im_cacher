package server

import (
	"im_cacher/src/aof"
	"im_cacher/src/cache"
	"net"
)

func StartListen(addr string, dict *cache.Dict, aof *aof.AOF) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println(err.Error())
		}

		go HandleConnection(conn, dict, aof)
	}

}
