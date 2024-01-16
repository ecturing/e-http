package main

import (
	"ews/ehttp"
	"io"
)

func main() {
	r := ehttp.NewRouter()
	ehttp.Server(r, "/server", server, ehttp.POST)
	ehttp.Confirm(":8080")
}

func server(rq *ehttp.E_Request, rp *ehttp.E_Response) {
	data, err := io.ReadAll(rq.Body)
	if err != nil {
		panic(err)
	}
	rp.DataFrom = string(data)
}
