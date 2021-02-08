package utils

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MarshalMessageBody convenience method to marshal different interfaces into JSON
func MarshalMessageBody(body interface{}) ([]byte, error) {
	switch msg := body.(type) {
	case proto.Message:
		marshaller := protojson.MarshalOptions{} // perhaps inject a custom one ?
		jsonBytes, err := marshaller.Marshal(msg)
		return jsonBytes, err
	case []byte:
		return msg, nil
	default:
		return json.Marshal(msg)
	}
}
