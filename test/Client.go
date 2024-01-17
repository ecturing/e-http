package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	GET  = "GET"
	POST = "POST"
	ERR  = "GTE"
	PUT  = "PUT"
)

func main() {
	test()
}

func test() {
	type args struct {
		method string
		addr   string
		body   string
	}
	tests := []struct {
		name string
		args args
		want string
		err  bool
	}{
		//生成测试用例，want参数为服务器返回参数，返回内容是客户端发送的body内容，利用参数实现最大覆盖测试,正确路由路径为server
		{
			name: "Test1",
			args: args{
				method: GET,
				addr:   "http://localhost:8080/server/v1/user",
				body:   "h",
			},
			want: "Hello World",
			err:  true,
		},
		{
			name: "Test2",
			args: args{
				method: POST,
				addr:   "http://localhost:8080/api/v1/user",
				body:   "id=1&name=2",
			},
			want: "id=1&name=2",
			err:  true,
		},
		{
			name: "Test3",
			args: args{
				method: POST,
				addr:   "http://localhost:8080/server",
				body:   "id=1&name=2",
			},
			want: "id=1&name=2",
			err:  false,
		},
	}

	for _, tt := range tests {
		if a := serverTest(tt.args.method, tt.args.addr, tt.args.body); a == tt.want {
			fmt.Printf("%s测试成功\n", tt.name)
		} else if tt.err {
			fmt.Printf("%s测试成功\n", tt.name)
		} else {
			fmt.Printf("%s测试失败\n", tt.name)
		}
	}
}

func serverTest(method string, addr string, body string) string {
	buff := bytes.NewBufferString(body)
	length := buff.Len()
	req, err := http.NewRequest(method, addr, buff)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	// 指定 HTTP 版本
	req.Proto = "HTTP/1.0"
	req.ProtoMajor = 1
	req.ProtoMinor = 0
	req.ContentLength = int64(length)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(bodyBytes)
}
