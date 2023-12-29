package main

import (
	"ews/ehttp"
	"testing"
)

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

func Test_server(t *testing.T) {
	type args struct {
		rq *ehttp.Request
		rp *ehttp.Response
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server(tt.args.rq, tt.args.rp)
		})
	}
}
