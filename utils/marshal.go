package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func MarshalMessageBody(body interface{}) ([]byte, error) {
	switch msg := body.(type) {
	case proto.Message:
		marshaller := jsonpb.Marshaler{} // perhaps inject a custom one ?
		jsonString, err := marshaller.MarshalToString(msg)
		return []byte(jsonString), err
	case []byte:
		return msg, nil
	default:
		return json.Marshal(msg)
	}
}
