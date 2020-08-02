package utils

import "google.golang.org/grpc/metadata"

type MDTraceCarrier metadata.MD

func (md MDTraceCarrier) Set(key, value string) {
	metadata.MD(md).Set(key, value)
}

func (md MDTraceCarrier) ForeachKey(handler func(key, value string) error) error {
	for k, vv := range md {
		for _, v := range vv {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

