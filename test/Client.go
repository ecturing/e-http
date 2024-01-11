package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	// "strconv"
)

// 写一个能向localhost:8080端口发送get,post请求的函数并编写测试用例，生成get请求，post请求各三个，请求参数分别为张三李四王五
func sendPostRequest() {
	resp, err := http.Post("http://localhost:8080/server", "application/json", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Body:", string(bodyBytes))
}

func sendGetRequest() {
	//请求体必须以双\r\n结尾
	buff := bytes.NewBufferString("张三")
	length := buff.Len()
	req, err := http.NewRequest("POST", "http://localhost:8080/server", buff)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 指定 HTTP 版本，例如 HTTP/2
	req.Proto = "HTTP/1.0"
	req.ProtoMajor = 1
	req.ProtoMinor = 0
	req.ContentLength = int64(length)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Body:", string(bodyBytes))
}

// 生成mian函数并调用上面的函数
func main() {
	sendGetRequest()
	// sendPostRequest()
}
