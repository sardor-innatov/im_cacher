package server

import (
	"fmt"
	"im_cacher/src/aof"
	"im_cacher/src/cache"
	"im_cacher/src/protocol"
	"im_cacher/src/protocol/resp"
	"net"
)

func HandleConnection(conn net.Conn, d *cache.Dict, aof *aof.AOF) {

	reader := resp.NewReader(conn)
	for {
		rawValue, _, err := reader.ReadValue()
		if err != nil {
			println(err.Error())
			break
		}

		request, rawRequest, err := protocol.SerializeRequest(rawValue)
		if err != nil {
			// error sendback logic
			println(err.Error())
			continue
		}

		switch request.Type {
		case "set":

			d.Set(request.Key, request.Value, request.TTL)
			aof.Append(*rawRequest)
			fmt.Fprint(conn, "ok")

		case "get":

			value, ok := d.Get(request.Key)
			if !ok {
				fmt.Fprint(conn, "nothing")
				continue
			}

			rawRESPbytes, err := value.MarshalRESP()
			//_, err := cacheItem.Value.MarshalRESP()
			if err != nil{
				println(err.Error())
			}
			conn.Write(rawRESPbytes)
			//fmt.Fprint(conn, cacheItem.Value)
			fmt.Println(rawRESPbytes)

		case "del":
			d.Del(request.Key)
			aof.Append(*rawRequest)
			fmt.Fprint(conn, "ok")
		}

	}
}
