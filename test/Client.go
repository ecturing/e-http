package test

import (
        "fmt"
        "io"
        "net/http"
)

//写一个能向localhost:8080端口发送get请求的函数
func sendGetRequest() {
    // 编写测试用例，生成get请求，post请求各三个，请求参数分别为张三李四王五
    



    resp, err := http.Get("http://localhost:8080/server")
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
//生成mian函数并调用上面的函数
func main() {
    sendGetRequest()
}