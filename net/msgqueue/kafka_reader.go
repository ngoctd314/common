package queue

import "io"

var _ io.ReadCloser = (*KafkaReader)(nil)

type KafkaReader struct {
}

func (k *KafkaReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// Close implements io.ReadCloser.
func (k *KafkaReader) Close() error {
	panic("unimplemented")
}
