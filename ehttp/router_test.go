package ehttp

import (
	"reflect"
	"testing"
)

func TestRouter_Register(t *testing.T) {
	type fields struct {
		root *treeNode
	}
	type args struct {
		pattern string
		f       ServerHTTP
		method  RequestMethod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test1", fields{root: NewRouter().root}, args{"/api/v1", func(r *E_Request, rp *E_Response) {}, GET}},
		{"test2", fields{root: NewRouter().root}, args{"/api/v2", func(r *E_Request, rp *E_Response) {}, GET}},
		{"test3", fields{root: NewRouter().root}, args{"/api/v3", func(r *E_Request, rp *E_Response) {}, POST}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				root: tt.fields.root,
			}
			r.Register(tt.args.pattern, tt.args.f, tt.args.method)
		})
	}
}

func TestRouter_Search(t *testing.T) {
	r:=NewRouter()
	r.Register("/api/v1", func(r *E_Request, rp *E_Response) {}, GET)
	r.Register("/api/v2", func(r *E_Request, rp *E_Response) {}, GET)
	r.Register("/api/v3", func(r *E_Request, rp *E_Response) {}, GET)
	type fields struct {
		root *treeNode
	}
	type args struct {
		pattern string
		method  RequestMethod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ServerHTTP
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", fields{root: r.root}, args{"/api/v1",GET}, func(r *E_Request, rp *E_Response) {}, false},
		{"test2", fields{root: r.root}, args{"/api/v5",GET}, nil , true},
		{"test3", fields{root: r.root}, args{"/api/v3",POST}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				root: tt.fields.root,
			}
			got, err := r.Search(tt.args.pattern, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Router.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Router.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
