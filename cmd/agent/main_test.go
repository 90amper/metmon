package main

// go build -o server *.go и go build -o agent *.go
// .\metricstest-windows-amd64.exe -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent
import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
