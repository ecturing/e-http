package ehttp

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

// 1.ehttp初步测试完成，2024/1/6

func Test_getPathValues(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name           string
		args           args
		wantRouterPath string
		wantValues     string
		wantErr        bool
	}{
		//生成多个测试用例，url为标准url格式 scheme://host:port/path?query，生成测试用例包含标准GET类请求url，标准POSTurl，不规范的url
		{
			name: "test1",
			args: args{
				url: "http://www.example.com/path/data?key=value",
			},
			wantRouterPath: "/path/data",
			wantValues:     "key=value",
			wantErr:        false,
		},
		{
			name: "test2",
			args: args{
				url: "https:/www.example.com/path/data?key=value&key2=value2",
			},
			wantRouterPath: "",
			wantValues:     "",
			wantErr:        true,
		},
		{
			name: "test3",
			args: args{
				url: "http://www.example.com/path/data?",
			},
			wantRouterPath: "/path/data",
			wantValues:     "",
			wantErr:        false,
		},
		{
			name: "test4",
			args: args{
				url: "http://www.example.com/path/data",
			},
			wantRouterPath: "/path/data",
			wantValues:     "",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRouterPath, gotValues, err := ReadPathValues(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPathValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRouterPath != tt.wantRouterPath {
				t.Errorf("getPathValues() gotRouterPath = %v, want %v", gotRouterPath, tt.wantRouterPath)
			}
			if gotValues != tt.wantValues {
				t.Errorf("getPathValues() gotValues = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

// 读取bufio.Reader里面的内容，返回string
func readBuffer(reader *bufio.Reader) string {
	var result strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		result.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
	return result.String()
}

func Test_getRequestLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantMethod string
		wantURL    string
		wantProto  string
		wantErr    bool
	}{
		// TODO: Add test cases.
		// 生成请求体第一行，包含请求方法、请求地址、协议类型，分别测试正常请求和不规范请求
		{"test1", args{"GET / HTTP/1.1"}, "GET", "/", "HTTP/1.1", false},
		{"test2", args{"POST / HTTP/1.1"}, "POST", "/", "HTTP/1.1", false},
		{"test3", args{"GET /path/data?key=value HTTP/1.1"}, "GET", "/path/data?key=value", "HTTP/1.1", false},
		{"test4", args{"POST /path/data?key=value HTTP/1.1"}, "POST", "/path/data?key=value", "HTTP/1.1", false},
		{"test5", args{"GET /path/data?key=value"}, "GET", "/path/data?key=value", "", true},
		{"test6", args{"POST /path/data?key=value"}, "POST", "/path/data?key=value", "", true},
		{"test7", args{"GET /path/data?key=value"}, "GET", "/path/data?key=value", "", true},
		{"test8", args{"POST /path/data?key=value"}, "POST", "/path/data?key=value", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotURL, gotProto, err := ReadRequestLine(tt.args.s)
			// 除了对这四个返回值进行测试，如果err!=nil,且tt.wantErr==true,则测试通过
			if (err != nil) && tt.wantErr == true {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getRequestLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("getRequestLine() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotURL != tt.wantURL {
				t.Errorf("getRequestLine() gotURL = %v, want %v", gotURL, tt.wantURL)
			}
			if gotProto != tt.wantProto {
				t.Errorf("getRequestLine() gotProto = %v, want %v", gotProto, tt.wantProto)
			}
		})
	}
}

func TestGetLineHeader(t *testing.T) {
	type args struct {
		b *bufio.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{"test1", args{bufio.NewReader(strings.NewReader("GET / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n\r\n"))}, "GET / HTTP/1.1", "Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n"},
		{"test2", args{bufio.NewReader(strings.NewReader("POST / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 18\r\n\r\n" +
			"{\"key\":\"value\"}"))}, "POST / HTTP/1.1", "Host: www.example.com\r\n" + "User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 18\r\n"},
		{"test3", args{bufio.NewReader(strings.NewReader("POST / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 13\r\n\r\n" +
			"key1=value1"))}, "POST / HTTP/1.1", "Host: www.example.com\r\n" + "User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 13\r\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLineHeader(tt.args.b)
			if got != tt.want {
				t.Errorf("GetLineHeader() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetLineHeader() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReadHeader(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.

		{"test1", args{"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n"}, map[string]string{"Host": "www.example.com", "User-Agent": "Go-http-client/1.1", "Accept": "*/*", "Accept-Encoding": "gzip"}},
		{"test2", args{"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n"}, map[string]string{"Host": "www.example.com", "User-Agent": "Go-http-client/1.1", "Accept": "*/*", "Accept-Encoding": "gzip"}},
		{"test3", args{"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n"}, map[string]string{"Host": "www.example.com", "User-Agent": "Go-http-client/1.1", "Accept": "*/*", "Accept-Encoding": "gzip"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadHeader(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBody(t *testing.T) {
	type args struct {
		reader *bufio.Reader
		length int
	}
	tests := []struct {
		name string
		args args
		want io.ReadCloser
	}{
		// TODO: Add test cases.
		{"test1", args{bufio.NewReader(strings.NewReader("GET / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip\r\n\r\n")), 0}, io.NopCloser(strings.NewReader(""))},
		{"test2", args{bufio.NewReader(strings.NewReader("POST / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 18\r\n\r\n" +
			"{\"key\":\"value\"}")), bytes.NewBufferString("{\"key\":\"value\"}").Len()}, io.NopCloser(strings.NewReader("{\"key\":\"value\"}"))},
		{"test3", args{bufio.NewReader(strings.NewReader("POST / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 13\r\n\r\n" +
			"key1=value1")), bytes.NewBufferString("key1=value1").Len()}, io.NopCloser(strings.NewReader("key1=value1"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := io.ReadAll(ReadBody(tt.args.reader, tt.args.length))
			want, _ := io.ReadAll(tt.want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ReadBody() = %v, want %v", got, want)
			}
		})
	}
}
