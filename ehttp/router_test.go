package ehttp

import (
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		want *Router
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// 生成3个测试用例，pattern分别为/api/v1依次v2v3，method为GET，f为ServerHTTP类型函数但函数体为空
		{"test1", fields{root: NewRouter().root}, args{"/api/v1", func(r *Request, rp *Response) {}, GET}},
		{"test2", fields{root: NewRouter().root}, args{"/api/v2", func(r *Request, rp *Response) {}, GET}},
		{"test3", fields{root: NewRouter().root}, args{"/api/v3", func(r *Request, rp *Response) {}, GET}},
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
	r.Register("/api/v1", func(r *Request, rp *Response) {}, GET)
	r.Register("/api/v2", func(r *Request, rp *Response) {}, GET)
	r.Register("/api/v3", func(r *Request, rp *Response) {}, GET)
	type fields struct {
		root *treeNode
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ServerHTTP
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", fields{root: r.root}, args{"/api/v1"}, func(r *Request, rp *Response) {}, false},
		{"test2", fields{root: r.root}, args{"/api/v5"}, nil , true},
		{"test3", fields{root: r.root}, args{"/api/v3"}, func(r *Request, rp *Response) {}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				root: tt.fields.root,
			}
			got, err := r.Search(tt.args.pattern)
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

func TestRouter_RouterListener(t *testing.T) {
	type fields struct {
		root *treeNode
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := &Router{
				root: tt.fields.root,
			}
			router.RouterListener()
		})
	}
}
