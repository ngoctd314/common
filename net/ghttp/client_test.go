package ghttp

import "testing"

func TestNewClient(t *testing.T) {
	t.Parallel()

	type args struct {
		opts []clientOption
	}
	testCases := []struct {
		name string
		args args
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			NewClient(tc.args.opts...)
		})
	}
}
