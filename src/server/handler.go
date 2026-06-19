package server

import (
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
		case "ping":
			conn.Write([]byte("+pong\r\n"))
		case "set":

			d.Set(request.Key, request.Value, request.TTL)
			aof.Append(*rawRequest)
			conn.Write([]byte("+ok\r\n"))

		case "get":

			value, ok := d.Get(request.Key)
			if !ok {
				nullResponse, _ := resp.NullValue().MarshalRESP()
				conn.Write(nullResponse)
				continue
			}

			rawRESPbytes, err := value.MarshalRESP()
			if err != nil {
				println(err.Error())
			}
			conn.Write(rawRESPbytes)

		case "del":
			d.Del(request.Key)
			aof.Append(*rawRequest)
			conn.Write([]byte("+ok\r\n"))
		}

	}
}
