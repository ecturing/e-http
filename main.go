package main

import (
	"ews/ehttp"
	"fmt"
)

func main() {
	ehttp.Server("/server",server)
	err:=ehttp.Confirm(":8080")
	if err!=nil {
		fmt.Println(err)
	}
}

func server(rq *ehttp.Request,rp *ehttp.Response){
	
}