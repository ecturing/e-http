package main

import (
	"ews/ehttp"
	"io"
)

func main() {
	r := ehttp.NewRouter()
	ehttp.ServerMux(r, "/server", server, ehttp.POST)
	ehttp.ListenAddr(":8080",r)
}

func server(rq *ehttp.E_Request, rp *ehttp.E_Response) {
	data, err := io.ReadAll(rq.Body)
	if err != nil {
		panic(err)
	}
	rp.DataFrom = string(data)
}
