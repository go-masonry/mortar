package utils

import "google.golang.org/grpc/metadata"

// MDTraceCarrier is an implementation for Tracing carrier, it will hold traceID, etc
type MDTraceCarrier metadata.MD

// Set part of the Carrier interface
func (md MDTraceCarrier) Set(key, value string) {
	metadata.MD(md).Set(key, value)
}

// ForeachKey part of the Carrier interface
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
