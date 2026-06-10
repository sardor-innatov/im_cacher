package protocol

import (
	"im_cacher/src/protocol/resp"
	"time"
)

type Request struct {
	Type  reqType
	TTL   time.Duration
	Key   string
	Value resp.Value
}

type reqType string

const (
	SET reqType = "set"
	GET reqType = "get"
	DEL reqType = "del"
)

func SerializeRequest(value resp.Value) (*Request, *resp.Value, error) {

	var request Request
	arrValue := value.Array()

	switch arrValue[0].String() {
	case "set":
		request.Type = SET
	case "get":
		request.Type = GET
	case "del":
		request.Type = DEL
	default:
		return nil, nil, ErrInvalidFormat
	}

	ttl := arrValue[1].Integer()
	if ttl <= 0 {
		request.TTL = 0
	} else if ttl > 0 {
		request.TTL = time.Second * time.Duration(ttl)
	}

	request.Key = arrValue[2].String()

	request.Value = arrValue[3]

	return &request, &value, nil

}
